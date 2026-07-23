# System Architecture Diagram — V4

```
╔══════════════════════════════════════════════════════════════════════════════════════════════╗
║                         JOB PROCESSING SYSTEM — V4 FINAL ARCHITECTURE                       ║
╚══════════════════════════════════════════════════════════════════════════════════════════════╝


  HTTP Client
      │
      │  POST /jobs   GET /jobs/{id}   GET /healthz
      ▼
╔═════════════════════════════════════════╗
║            API SERVICE                  ║
║        cmd/api/main.go · :8080          ║
║                                         ║
║  POST /jobs:                            ║
║  ┌─────────────────────────────────┐    ║
║  │ BEGIN TRANSACTION                │    ║
║  │   INSERT INTO jobs (pending)     │    ║
║  │   INSERT INTO outbox (job.created│    ║──────────────────────────────────────────────┐
║  │ COMMIT → 202 Accepted            │    ║  (atomic write — dual-write                  │
║  └─────────────────────────────────┘    ║   problem solved)                            │
║                                         ║                                               │
║  GET /jobs/{id}:                        ║                                               │
║    SELECT FROM jobs → 200/404           ║◄──────── SELECT jobs ─────────────────────┐  │
╚═════════════════════════════════════════╝                                            │  │
                                                                                       │  │
                                                                                       │  │
╔═════════════════════════════════════════╗                                            │  │
║          PUBLISHER SERVICE              ║                                            │  │
║          cmd/publisher/main.go          ║                                            │  │
║                                         ║                                            │  │
║  loop every 1s:                         ║                                            │  │
║    BEGIN TX                             ║                                            │  │
║    SELECT outbox WHERE published_at     ║◄──────── SELECT outbox ─────────────┐     │  │
║      IS NULL FOR UPDATE SKIP LOCKED     ║          (SKIP LOCKED, LIMIT 100)   │     │  │
║                                         ║                                      │     │  │
║    for each event:                      ║                                      │     │  │
║      LPUSH jobID    → jobs:pending  ────╬════════════════════════════╗         │     │  │
║      LPUSH envelope → events:pending ───╬════════════════════════╗   ║         │     │  │
║                                         ║                        ║   ║         │     │  │
║    UPDATE outbox SET published_at=now() ║────────────────────────╬───╬─────────┤     │  │
║    COMMIT                               ║                        ║   ║         │     │  │
╚═════════════════════════════════════════╝                        ║   ║         │     │  │
                                                                   ║   ║         │     │  │
                                          ┌────────────────────────╫───╫─────────┤     │  │
                                          │         PostgreSQL      ║   ║         │     │  │
                                          │                         ║   ║         │     │  │
                                          │  jobs                   ║   ║◄────────┘     │  │
                                          │    id, payload          ║   ║               │  │
                                          │    status               ║   ║               │  │
                                          │    result, error        ║   ║               │  │
                                          │    retry_count          ║   ║◄──────────────┘  │
                                          │    retry_after          ║   ║                  │
                                          │                         ║   ║◄─────────────────┘
                                          │  outbox                 ║   ║
                                          │    id, event_type       ║   ║
                                          │    payload (job_id)     ║   ║
                                          │    created_at           ║   ║
                                          │    published_at         ║   ║
                                          │                         ║   ║
                                          │  processed_events       ║   ║
                                          │    event_id (PK)        ║   ║◄──── INSERT (consumer)
                                          │    processed_at         ║   ║
                                          └─────────────────────────╫───╫───────────────────┐
                                                                     ║   ║                   │
                                          ┌──────────────────────────╫───╫──────────────┐   │
                                          │         Redis             ║   ║              │   │
                                          │                           ║   ║              │   │
                                          │  jobs:pending  ◄══════════╝   ║              │   │
                                          │  (LPUSH by publisher)         ║              │   │
                                          │  (BLMOVE by dispatcher)       ║              │   │
                                          │                               ║              │   │
                                          │  jobs:processing              ║              │   │
                                          │  (BLMOVE target — atomic      ║              │   │
                                          │   claim, LREM on done)        ║              │   │
                                          │                               ║              │   │
                                          │  events:pending ◄═════════════╝              │   │
                                          │  (LPUSH by publisher)                        │   │
                                          │  (BLPOP by consumer)                         │   │
                                          └──────────────────────┬───────────────────────┘   │
                                                                  │                           │
                        ┌─────────────────────────────────────────┘                           │
                        │  BLMOVE jobs:pending → jobs:processing                              │
                        │                                                                     │
                        ▼                                                                     │
╔═════════════════════════════════════════╗                                                   │
║           WORKER SERVICE                ║                                                   │
║           cmd/worker/main.go            ║                                                   │
║                                         ║                                                   │
║  Dispatcher:                            ║                                                   │
║    BLMOVE jobs:pending                  ║                                                   │
║         → jobs:processing (atomic)      ║                                                   │
║    send jobID to worker channel         ║                                                   │
║                                         ║                                                   │
║  Worker goroutines (×N, default 5):     ║                                                   │
║    UPDATE jobs SET status=processing    ║──────────────────────────────────────────────┐   │
║    run AI agentic loop ────────────────►║  Anthropic API                               │   │
║      (tool calls, multi-step)           ║                                               │   │
║    on success:                          ║                                               │   │
║      UPDATE jobs SET status=completed   ║──────────────────────────────────────────────┤   │
║      LREM jobs:processing               ║──────────────────────────── Redis            │   │
║    on failure:                          ║                                               │   │
║      retry_count++, backoff 2^n sec     ║                                               │   │
║      UPDATE jobs SET status=pending     ║──────────────────────────────────────────────┤   │
║      LREM jobs:processing               ║──────────────────────────── Redis            │   │
║                                         ║                                               │   │
║  Reaper (every 60s):                    ║                                               │   │
║    SELECT processing WHERE              ║                                               │   │
║      updated_at < now() - 300s          ║──────────────────────────────────────────────┤   │
║    UPDATE status=pending, retry_count++ ║   all job state reads/writes → PostgreSQL    │   │
║    LREM jobs:processing                 ║──────────────────────────── Redis            │   │
╚═════════════════════════════════════════╝                                               │   │
                                                                                          │   │
                        ┌────────────────────── BLPOP events:pending ─────────────────────────┘
                        │
                        ▼
╔═════════════════════════════════════════╗
║          CONSUMER SERVICE               ║
║          cmd/consumer/main.go           ║
║                                         ║
║  loop:                                  ║
║    BLPOP events:pending (5s timeout)    ║
║    decode EventEnvelope{                ║
║      event_id, event_type, job_id}      ║
║                                         ║
║    INSERT INTO processed_events         ║──────────────────────────────────────────────┐
║    (event_id) ON CONFLICT → skip        ║   idempotency: PK constraint enforces        │
║                                         ║   at-most-once processing                    │
║    if inserted: log event processed     ║──────────────────────────── PostgreSQL       ◄┘
║    if conflict: log duplicate, skip     ║
╚═════════════════════════════════════════╝


══════════════════════════════════════ GRACEFUL SHUTDOWN ══════════════════════════════════════

  SIGTERM / SIGINT
        │
        │  signal.NotifyContext() cancels shared ctx in each process
        ▼
  ┌──────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────┐
  │   API            │  │  PUBLISHER          │  │  CONSUMER           │  │  WORKER         │
  │                  │  │                     │  │                     │  │                 │
  │  <-ctx.Done()    │  │  <-ctx.Done()       │  │  BLPOP returns      │  │  <-ctx.Done()   │
  │  srv.Shutdown()  │  │  poll loop exits    │  │  context error      │  │  close channel  │
  │  drain in-flight │  │  current batch      │  │  loop returns nil   │  │  workers finish │
  │  HTTP requests   │  │  commits or rolls   │  │                     │  │  current job    │
  │                  │  │  back cleanly       │  │                     │  │  wg.Wait()      │
  └──────────────────┘  └─────────────────────┘  └─────────────────────┘  └─────────────────┘


══════════════════════════════════════ DATA FLOW SUMMARY ══════════════════════════════════════

  1. JOB CREATION (atomic)
     Client → POST /jobs → API → [BEGIN TX: INSERT jobs + INSERT outbox] → COMMIT

  2. EVENT RELAY (at-least-once)
     Publisher → SELECT outbox (SKIP LOCKED) → LPUSH jobs:pending + LPUSH events:pending
               → UPDATE outbox.published_at → COMMIT
     (crash before UPDATE → event republished next iteration)

  3. JOB PROCESSING
     Dispatcher → BLMOVE jobs:pending → jobs:processing (atomic claim, no duplicate work)
     Worker     → AI loop → UPDATE jobs (completed/failed) → LREM jobs:processing

  4. STUCK JOB RECOVERY
     Reaper → SELECT processing > 300s → reset to pending → LREM jobs:processing

  5. EVENT CONSUMPTION (idempotent)
     Consumer → BLPOP events:pending → INSERT processed_events ON CONFLICT SKIP
     (duplicate event → PK violation → skip, log)

  6. JOB STATUS POLLING
     Client → GET /jobs/{id} → SELECT FROM jobs → 200 {status, result, error}
```
