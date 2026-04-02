# Architecture: Job Processing System

## What this system is

A backend service that accepts jobs via HTTP, processes them asynchronously, tracks their lifecycle, and returns results. It models the same fundamental problem behind task queues like Celery or Sidekiq: reliable, exactly-once, asynchronous job execution — even when workers crash or race.

The core constraint that makes this non-trivial: a job must be processed **reliably, exactly once, asynchronously**, under concurrent workers and partial failure conditions.

---

## V1 scope

V1 is minimal but correct. It solves the hard problems from the start — no shortcuts on locking or state transitions. Complexity is added in layers on top of this foundation.

---

## Components

### API layer
The only public surface. Accepts jobs, exposes status. Has no knowledge of workers or processing logic.

Endpoints:
- `POST /jobs` — validate payload, INSERT row with status=`pending`, return `{ id, status }`
- `GET /jobs/{id}` — read current job state from DB, return full record
- `GET /healthz` — liveness check

The API does not trigger workers. It does not wait for processing. It writes to the DB and returns immediately. The async contract is explicit: client gets an ID, polls for result.

### Persistence layer (Postgres)
Single source of truth for all job state.

`jobs` table:
```sql
id          UUID PRIMARY KEY DEFAULT gen_random_uuid()
status      TEXT NOT NULL DEFAULT 'pending'   -- pending | processing | completed | failed
payload     JSONB NOT NULL
result      JSONB
error       TEXT
retry_count INT NOT NULL DEFAULT 0
retry_after TIMESTAMPTZ
created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
```

### Worker
Background goroutine. No HTTP surface. Loops continuously:
1. Poll DB for a claimable job
2. Lock it atomically
3. Execute processing logic
4. Write result or error back

**Critical:** polling uses `SELECT FOR UPDATE SKIP LOCKED` — this is non-negotiable. A naive `SELECT` creates a race condition the moment two workers run. `SKIP LOCKED` atomically skips rows already locked by another worker, preventing duplicate processing without application-level coordination.

Poll query:
```sql
SELECT * FROM jobs
WHERE status = 'pending'
  AND (retry_after IS NULL OR retry_after < NOW())
ORDER BY created_at
LIMIT 1
FOR UPDATE SKIP LOCKED
```

The status update to `processing` happens **inside the same transaction** as the SELECT. This is atomic — there is no window between claiming and locking.

### State machine
Not a separate service. A set of rules enforced inside DB transactions.

```
pending → processing → completed
                     ↘ failed (if retry_count >= max_retries)
         ↗ (retry)
```

Valid transitions only. A job never moves backwards arbitrarily. The only backward move is `processing → pending` on retry, and only when `retry_count < max_retries`.

### Reliability mechanisms

**Locking:** `SELECT FOR UPDATE SKIP LOCKED` — prevents two workers from claiming the same job.

**Retries:** On worker failure, increment `retry_count`, set `retry_after = now() + backoff`, reset status to `pending`. If `retry_count >= max_retries`, set status to `failed`.

**Backoff:** `retry_after` column prevents immediate retry hammering. Start with exponential backoff (e.g. 2^retry_count seconds).

**Failure capture:** `error` column stores the last error message. Always persisted before marking `failed`.

**Stuck job recovery (V2):** Jobs stuck in `processing` (worker crashed mid-job) need a timeout check — a separate process resets jobs where `updated_at < now() - interval '5 minutes'` and status is still `processing`.

---

## Project structure

```
cmd/
  server/
    main.go               ← wire everything, start HTTP + worker
internal/
  api/
    handler.go            ← HTTP handlers (no business logic, no SQL)
    routes.go             ← route registration
    server.go             ← http.Server setup
  worker/
    worker.go             ← poll loop, job execution
    dispatcher.go         ← goroutine pool (V2)
  job/
    job.go                ← Job type, status constants, state machine rules
    store.go              ← DB queries (all SQL lives here)
  db/
    db.go                 ← connection setup
migrations/
  001_create_jobs.sql
```

**Boundary rules:**
- Handlers call store, never raw SQL
- Workers call store, never HTTP
- Store returns domain types (`job.Job`), not raw DB rows
- No business logic in handlers or store — only in `job/`

---

## Architecture decisions

### Why polling over event-driven (V1)
Polling forces explicit ownership of the hard problems: race conditions, locking, atomic state transitions. With a queue, the queue absorbs that complexity. Polling makes it visible and teachable. Event-driven is V2 — after polling friction is felt.

### Why SELECT FOR UPDATE SKIP LOCKED
Standard `SELECT WHERE status='pending'` + separate `UPDATE` is two operations. Any other worker can claim the same row between them. `SKIP LOCKED` makes claim + lock a single atomic operation at the DB level. No application-level mutex needed.

### Why retry_after instead of immediate retry
Retrying immediately after failure hammers the system and obscures real failure rates. `retry_after` with exponential backoff gives downstream systems time to recover and makes retry behavior observable.

### Why 202 vs 201 on POST /jobs
`202 Accepted` is more honest — the resource exists but processing hasn't started. `201 Created` implies the work is done. Use `202` to make the async contract explicit to clients.

---

## Build order

1. `POST /jobs` + `GET /jobs/{id}` with DB writes (no worker yet)
2. Single worker with `SELECT FOR UPDATE SKIP LOCKED`
3. State transitions inside transactions
4. Retry logic with `retry_count` + `retry_after`
5. Dispatcher (goroutine pool) — V2
6. Stuck job recovery — V2
7. Queue (Redis/Kafka) replacing polling — V3

Each step is a discrete, testable addition. V1 = steps 1–4.

---

## What this is NOT (V1)
- No external queue (Redis/Kafka) — optional V2+
- No outbox pattern — optional V2+
- No metrics/tracing — optional V2+
- No multiple service instances — optional V2+
- No complex business logic in job execution
- No UI