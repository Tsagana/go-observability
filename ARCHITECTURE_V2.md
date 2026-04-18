# Architecture: Job Processing System — V2

## What V2 is

V2 takes V1's correct but serial system and makes it concurrent, resilient to worker crashes, and observable. The DB schema and state machine do not change. Everything is additive on top of V1.

V1 gave you: one worker, one job at a time, correct locking, retries, backoff.
V2 gives you: a worker pool, crash recovery, structured logging, and real AI processing.

---

## V2 components

### V2a — Dispatcher + worker pool

**Problem it solves:**
V1 processes one job at a time. Naively spinning up N worker goroutines causes thundering herd — all workers poll the DB simultaneously, most find nothing. The dispatcher eliminates this.

**Design:**
One dispatcher owns the pool. Only the dispatcher polls the DB. Workers are dumb — they wait on a channel and process whatever arrives.

```
Dispatcher
  ├── poll loop         (only component that claims jobs from DB)
  ├── job channel       (buffered — acts as backpressure)
  └── worker pool
        worker 1 ─┐
        worker 2  ├── read from channel → execute → write result
        worker 3  │
        worker N ─┘
```

**Configuration:**
- `WORKER_COUNT` — number of worker goroutines. Configurable via env var, not hardcoded.
- `JOB_CHANNEL_BUFFER` — buffered channel size. Start with `WORKER_COUNT`. Too large = claimed jobs lost on crash. Too small = workers starve.
- `POLL_INTERVAL` — how often the dispatcher polls when no jobs are found. Default 2s.

**Backpressure:**
When the channel is full, the poller blocks naturally. The DB sees reduced polling load automatically. No explicit throttling needed.

**Graceful shutdown:**
Shutdown sequence on SIGTERM:
1. Cancel context — poller stops, no new jobs claimed
2. Close job channel — workers finish current job then exit
3. Wait for all workers via WaitGroup — no in-flight jobs killed
4. Process exits cleanly

```
SIGTERM
  → context cancelled
  → poller exits
  → job channel closed
  → workers drain current job, exit
  → WaitGroup reaches zero
  → process exits
```

Without graceful shutdown, a deploy kills workers mid-job. Those jobs get stuck in `processing` until the reaper recovers them.

---

### V2b — Stuck job reaper

**Problem it solves:**
A worker can die mid-job — panic, OOM kill, hard shutdown. The job stays in `processing` indefinitely. Graceful shutdown helps but does not cover hard crashes.

**Design:**
A background goroutine running on a timer. Every 60 seconds it scans for jobs stuck in `processing` beyond a configurable threshold and returns them to the queue.

**Recovery query:**
```sql
UPDATE jobs
SET
    status      = 'pending',
    retry_count = retry_count + 1,
    retry_after = now() + (interval '1 second' * power(2, retry_count)),
    updated_at  = now()
WHERE status     = 'processing'
  AND updated_at < now() - interval '5 minutes'
  AND retry_count < 3
RETURNING id;
```

Jobs exceeding `max_retries` are marked `failed` with an error message instead of returned to `pending`.

**Configuration:**
- `STUCK_JOB_TIMEOUT` — how long a job can be in `processing` before considered stuck. Must be longer than your p99 job duration. Default: 5 minutes.
- `REAPER_INTERVAL` — how often the reaper runs. Default: 60 seconds.

**Important interaction with V2d (AI):**
If AI API calls take up to 30 seconds, set `STUCK_JOB_TIMEOUT` to at least 2–3x that. False positives (reaping a still-running job) cause duplicate processing.

---

### V2c — Structured logging

**Problem it solves:**
With 5 concurrent workers processing jobs simultaneously, you cannot follow execution mentally. Structured logging makes the system observable and debuggable.

**Design:**
Use Go's `slog` package (stdlib since 1.21). Emit key-value structured log lines on every state transition. Every log line carries `job_id` and `worker_id` so you can trace a single job across concurrent output.

**Log events:**

| Event | Fields |
|---|---|
| job claimed | job_id, worker_id, previous_status |
| job completed | job_id, worker_id, duration_ms |
| job failed permanently | job_id, worker_id, error, retry_count |
| job retrying | job_id, retry_count, retry_after |
| job stuck recovered | job_id, stuck_duration |
| poll result | jobs_found, workers_idle |
| worker started | worker_id |
| worker stopped | worker_id |
| dispatcher shutdown | jobs_in_flight |

**Example log lines:**
```
level=INFO  msg="job.claimed"    job_id=abc-123 worker_id=2 retry_count=0
level=INFO  msg="job.completed"  job_id=abc-123 worker_id=2 duration_ms=1821
level=WARN  msg="job.retrying"   job_id=def-456 worker_id=1 error="rate limit" retry_after=30s
level=ERROR msg="job.failed"     job_id=ghi-789 worker_id=3 error="max retries exceeded"
level=INFO  msg="reaper.recover" job_id=jkl-012 stuck_duration=7m32s
```

**Why V2c before V2d:**
You need logging in place before integrating AI calls. Concurrent API calls without logging are nearly impossible to debug.

---

### V2d — AI processing

**Problem it solves:**
The worker stub (simulated sleep) doesn't test reliability patterns with real stakes. An AI API call does — it's slow, fails transiently, costs money on duplicates, and has rate limits that make backoff matter.

**Design:**
Replace the execute stub with a real Anthropic API call. One job type, one prompt, one result. Keep it simple — the architecture is the point, not the AI capability.

```
job payload:  { "text": "..." }
worker calls: Anthropic Claude API (summarization prompt)
job result:   { "summary": "..." }
```

**Error classification:**
Not all errors are equal. The worker must distinguish:

| Error type | Example | Action |
|---|---|---|
| Retryable | 429 Rate Limited, 503 Unavailable | retry with backoff |
| Permanent | 400 Bad Request, invalid prompt | mark failed immediately |
| Timeout | context deadline exceeded | retry with backoff |

**The idempotency problem:**
Sequence that causes duplicate API calls:
1. Worker calls Claude API — succeeds, result returned
2. Worker begins writing result to DB
3. Worker crashes before commit
4. Job recovered by reaper, retried
5. Worker calls Claude API again — paid twice, different result possible

The fix: write the API response and mark the job complete in a single transaction. No gap between "got result" and "stored result."

**Timeout handling:**
Set a context timeout on the API call. Must be shorter than `STUCK_JOB_TIMEOUT` so the worker times out and retries cleanly before the reaper considers the job stuck.

**Configuration:**
- `ANTHROPIC_API_KEY` — API key, never hardcoded
- `AI_TIMEOUT` — per-call timeout. Default: 30s
- `AI_MAX_RETRIES` — max retry attempts for rate limits. Separate from job max_retries.

---

## Updated project structure

```
cmd/
  server/
    main.go                   ← wire dispatcher, reaper, API server, handle SIGTERM

internal/
  api/
    handler.go                ← HTTP handlers (unchanged from V1)
    routes.go
    server.go

  worker/
    dispatcher.go             ← NEW: owns pool, poll loop, channel, WaitGroup
    worker.go                 ← updated: reads from channel instead of polling
    reaper.go                 ← NEW: stuck job recovery loop
    executor.go               ← NEW: job execution logic (AI call lives here)

  job/
    job.go                    ← Job type, status constants (unchanged)
    store.go                  ← DB queries (add reaper query)
    errors.go                 ← NEW: RetryableError type, error classification

  ai/
    client.go                 ← NEW: Anthropic API client wrapper
    prompt.go                 ← NEW: prompt construction

  db/
    db.go                     ← unchanged

migrations/
  001_create_jobs.sql         ← unchanged (schema does not change in V2)
```

---

## Configuration reference

All configuration via environment variables:

```
# Worker pool
WORKER_COUNT=5
JOB_CHANNEL_BUFFER=5
POLL_INTERVAL=2s

# Reaper
STUCK_JOB_TIMEOUT=5m
REAPER_INTERVAL=60s

# Jobs
MAX_RETRIES=3

# AI (V2d)
ANTHROPIC_API_KEY=sk-...
AI_TIMEOUT=30s

# DB (unchanged from V1)
DATABASE_URL=postgres://...
```

---

## Architecture decisions

### Why one poller, not one poller per worker
Multiple pollers cause thundering herd — all hitting the DB at the same interval, creating spikes of load and lock contention. One poller creates one steady stream of claimed jobs fed to workers via channel. The channel acts as a natural buffer and decouples polling rate from worker speed.

### Why a buffered channel sized to worker count
Too small: workers go idle waiting for the poller to claim the next job — throughput drops.
Too large: many jobs are claimed and held in memory. If the process crashes, those jobs are stuck in `processing` until the reaper recovers them. Sizing to `WORKER_COUNT` means at most N extra jobs are held in memory at any time.

### Why graceful shutdown before AI integration
An in-flight AI call takes up to 30 seconds. Without graceful shutdown, a SIGTERM kills it instantly — job stuck, API call paid for, result lost. With graceful shutdown, the worker finishes the call, writes the result, then exits. Order matters.

### Why classify errors before retrying
Retrying a permanent error wastes retry budget and delays the job reaching `failed` state. A `400 Bad Request` on a malformed prompt will never succeed — retrying it 3 times just delays the inevitable and consumes resources. Classify on first failure, route accordingly.

### Why AI timeout must be less than stuck job threshold
If `AI_TIMEOUT=30s` and `STUCK_JOB_TIMEOUT=5m`, the worker will time out and retry the job cleanly long before the reaper considers it stuck. If the timeout were longer than the threshold, the reaper would recover a job that is still being actively processed — causing duplicate processing.

---

## Build order

Build in this exact sequence. Each step is testable before moving to the next.

```
V2a  Dispatcher + worker pool + graceful shutdown
       ↓ verify: N workers process jobs concurrently, SIGTERM drains cleanly

V2b  Stuck job reaper
       ↓ verify: manually set a job to 'processing', confirm reaper recovers it

V2c  Structured logging
       ↓ verify: can trace a single job_id across concurrent worker output

V2d  AI integration
       ↓ verify: job payload processed by Claude, result stored, rate limit triggers retry
```

Do not implement V2d before V2c. Debugging concurrent AI calls without structured logging is extremely painful.

---

## What does not change from V1

- DB schema — no new columns, no new tables
- State machine — same transitions: pending → processing → completed / failed
- Locking mechanism — `SELECT FOR UPDATE SKIP LOCKED` unchanged
- Retry logic — `retry_count` + `retry_after` + exponential backoff unchanged
- API endpoints — `POST /jobs`, `GET /jobs/{id}`, `GET /healthz` unchanged
- Boundary rules — handlers never touch SQL, workers never touch HTTP

V2 is purely additive. If V1 is correct, V2 does not break it.
