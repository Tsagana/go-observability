# Architecture: Job Processing System — V4

## What V4 is

Two structural changes built on top of V3:

1. **Service split** — separate the API and the worker into independent processes
2. **Outbox pattern + consumer** — make event publishing atomic with the DB write, and add a consumer that processes those events idempotently

V4 turns a single-process system into a small distributed system. It introduces the patterns most cited in distributed systems interviews: transactional outbox, at-least-once delivery, and idempotent consumers.

---

## What does not change

- Postgres remains source of truth for job state
- Redis remains the job delivery layer (from V3)
- State machine, retry logic, AI loop, reaper logic — unchanged
- API endpoints and contracts — unchanged
- The `jobs` table schema is unchanged

V4 is additive in terms of components and changes only the wiring between them.

---

## Part 1: Service split

### Why
Up to V3, API and worker ran in one process. They share a DB and Redis but no in-process state. Splitting them lets each scale and deploy independently, and forces all coordination through infrastructure rather than function calls. This is the foundational move from monolith to multi-service architecture.

### New process layout

```
cmd/
  api/main.go         HTTP server only
  worker/main.go      dispatcher + worker pool + reaper
  publisher/main.go   outbox publisher (Part 2)
  consumer/main.go    event consumer (Part 2)
```

Each is a separate binary. Each has its own Docker container. All four share:

```
internal/
  job/         Job type, store, errors
  queue/       Queue interface, Redis + memory impls
  db/          Postgres connection
  ai/          Anthropic client, agent loop
  outbox/      Outbox table access, publisher logic (Part 2)
  events/      Event types, consumer logic (Part 2)
```

### Coordination
The services do not call each other. They coordinate via:

- **Postgres** — job state, outbox events, processed events
- **Redis** — job delivery queue, event delivery queue

The API never knows the worker exists. The worker never knows the API exists. They share infrastructure.

### What this teaches
- Service boundaries enforced by process boundaries
- Independent deployment and scaling
- Why "shared DB" between services has costs (schema coupling, contention)
- The conceptual difference between in-process coordination and distributed coordination

---

## Part 2: Outbox pattern

### The problem it solves
In V3, `POST /jobs` does two writes:
1. INSERT into Postgres `jobs` table
2. LPUSH job ID to Redis pending list

These are not atomic. If step 2 fails after step 1 succeeds, the job is orphaned in Postgres and the reaper recovers it eventually. That works but it's hand-wavy. There's a window during which the system is in an inconsistent state.

The outbox pattern eliminates that window.

### The pattern
Instead of writing directly to Redis, the API writes an event to an `outbox` table in the same transaction as the job insert. A separate publisher process reads the outbox and forwards events to Redis.

```
POST /jobs (single transaction):
  INSERT INTO jobs (...) VALUES (...)
  INSERT INTO outbox (event_type, payload) VALUES ('job.created', {jobID})
COMMIT

Publisher (separate process):
  loop:
    SELECT * FROM outbox WHERE published_at IS NULL ORDER BY id LIMIT 100
    For each:
      LPUSH to Redis events list
      UPDATE outbox SET published_at = now() WHERE id = $1
```

The DB transaction guarantees both writes succeed or neither does. The publisher's job is to bridge the DB to Redis. If the publisher crashes mid-batch, the unpublished events are still in the DB waiting for the next iteration.

### Outbox table schema

```sql
CREATE TABLE outbox (
    id           BIGSERIAL PRIMARY KEY,
    event_type   TEXT NOT NULL,
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);

CREATE INDEX outbox_unpublished_idx
    ON outbox (id)
    WHERE published_at IS NULL;
```

The partial index makes "find unpublished events" fast even when the table grows large.

### Publisher details

The publisher runs as a separate process (`cmd/publisher/main.go`). Its loop:

1. Begin transaction
2. SELECT unpublished events, FOR UPDATE SKIP LOCKED, LIMIT 100
3. For each event, LPUSH to Redis events list
4. UPDATE published_at for each successfully pushed event
5. COMMIT

`SKIP LOCKED` lets multiple publisher instances run safely if needed for scale. Most deployments only need one publisher.

If the publisher pushes to Redis but crashes before the UPDATE, the same event will be republished on the next iteration. This is at-least-once delivery — the consumer must handle duplicates.

### Why the API no longer pushes to Redis directly
After V4, the API only writes to Postgres. The job creation flow becomes:

```
POST /jobs:
  1. Validate payload
  2. Begin transaction
  3. INSERT into jobs (status = pending)
  4. INSERT into outbox ('job.created', {jobID})
  5. Commit
  6. Return 202
```

The publisher does the Redis push. This means the dispatcher's `BLMOVE` on Redis is now fed entirely by the publisher, not the API.

### Why keep published events
Don't delete them immediately. Keep the outbox as an event audit log. Add a separate cleanup job (cron or scheduled) that deletes events older than N days (e.g., 30 days). Storage is cheap; the audit trail is valuable when debugging "did we ever publish event X?"

---

## Part 3: Consumer service

### Why add a consumer
Without one, the outbox publishes events that nothing receives. The pattern is incomplete as a learning exercise. A consumer:

- Validates the at-least-once delivery contract end-to-end
- Forces you to implement idempotent consumption
- Closes the architectural loop: producer → outbox → publisher → queue → consumer

### Keep it tight
The consumer doesn't do anything meaningful. It logs events and tracks which it has seen. The point is the pattern, not the consumer's business logic.

### Consumer flow

```
loop:
  event = BLPOP from Redis events list (blocking)
  if event.id IN processed_events:
    log "duplicate, skipping"
    continue
  process event (log it)
  INSERT event.id INTO processed_events
```

### Processed events table

```sql
CREATE TABLE processed_events (
    event_id     BIGINT PRIMARY KEY,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Primary key constraint enforces idempotency. If the consumer crashes after processing but before inserting, the next attempt INSERT will succeed (no duplicate). If the consumer crashes after inserting but before completing processing... well, in this minimal version that's fine because "processing" is just logging.

For a real consumer, you'd want processing and the INSERT to happen in the same transaction. Worth noting in the design discussion even if not implemented.

---

## Updated project structure

```
cmd/
  api/
    main.go              HTTP server
  worker/
    main.go              dispatcher + pool + reaper
  publisher/
    main.go              outbox → Redis
  consumer/
    main.go              Redis → processed_events

internal/
  api/                   handlers (unchanged endpoints)
  worker/                dispatcher, worker, reaper, executor
  job/                   Job type, store, errors
  ai/                    Anthropic client, agent loop
  queue/                 Queue interface, Redis + memory
  db/                    Postgres connection
  redis/                 Redis connection
  outbox/                NEW: outbox table access, publisher logic
  events/                NEW: event types, consumer logic

migrations/
  001_create_jobs.sql            (V1)
  002_create_outbox.sql          (V4)
  003_create_processed_events.sql (V4)
```

---

## Updated configuration

New variables:

```
# Publisher
PUBLISHER_BATCH_SIZE=100
PUBLISHER_POLL_INTERVAL=1s
EVENTS_QUEUE_KEY=events:pending

# Consumer
CONSUMER_BLOCK_TIMEOUT=5s
```

All V1-V3 config remains unchanged.

---

## Failure modes and recovery

| Failure | Behavior |
|---|---|
| API crashes mid-request | Transaction not committed, no job created, no outbox event |
| API commits, publisher hasn't run yet | Event sits in outbox, picked up on next publisher iteration |
| Publisher pushes to Redis, crashes before UPDATE | Event republished next iteration, consumer sees duplicate and skips |
| Publisher crashes while running | Restart, resumes from oldest unpublished event |
| Consumer crashes after BLPOP, before processing | Event lost from Redis but publisher already marked published → permanent loss in this simple design. Real consumers use ACKs or transactional read-and-process. Note as known limitation. |
| Consumer crashes after processing, before INSERT to processed_events | Next event seen as duplicate, skipped — safe |
| Redis crashes, loses events list | Events still in outbox, publisher republishes everything not yet marked published |
| Postgres crashes | All services degrade gracefully, recover when DB returns |

The known limitation in the consumer crash case is intentional. A production consumer would do BRPOPLPUSH (move to a processing list) and only LREM after processing, matching the pattern the worker uses for jobs. Worth implementing if you want full coverage, but the basic version is enough to demonstrate the idempotency concept.

---

## Architecture decisions

### Why a separate publisher process
The API should be a thin write layer. Coupling Redis publishing into the API request handler ties HTTP latency to Redis availability. A separate publisher decouples the two and lets you scale them independently. It also keeps the API's job small: write to Postgres, return 202.

### Why not use Postgres LISTEN/NOTIFY instead of polling the outbox
LISTEN/NOTIFY works but has limitations: notifications are lost if no consumer is connected, payloads are limited in size, and it doesn't integrate cleanly with transactional outbox semantics. Polling the outbox table is simpler, durable, and the standard approach in production systems.

### Why one outbox table for all event types
Simplicity. The `event_type` column distinguishes them. If event volume grows enough to warrant separate tables (different retention, different consumers), splitting can come later. Premature partitioning is its own form of over-engineering.

### Why the consumer is dumb
The architectural point is the pattern — producer, outbox, publisher, broker, idempotent consumer. The consumer doing real work would obscure the pattern with business logic. A logging consumer makes the data flow visible without distraction.

### Why processed_events is a separate table, not a column on outbox
The consumer is a different service than the publisher. They shouldn't share a write surface. The consumer owns `processed_events`; the publisher owns `outbox.published_at`. Separation of concerns at the table level.

### Why no acknowledgement protocol between consumer and publisher
That would be a custom messaging protocol on top of Redis. Not worth it for V4. In a real system you'd use a broker that supports ACKs natively (RabbitMQ, Kafka with offset commits). The known limitation in the failure modes table is the honest acknowledgement of this gap.

---

## Build order

```
Weekend 1 — Service split
  - Create cmd/api/main.go (extract HTTP server)
  - Create cmd/worker/main.go (extract dispatcher + pool + reaper)
  - Update docker-compose to run both as separate services
  - Verify: kill one service, the other continues working
  - Verify: jobs submitted to API are processed by worker via DB + Redis

Weekend 2 — Outbox + publisher
  - Migration 002: create outbox table + partial index
  - Update POST /jobs handler: single transaction (jobs + outbox), no Redis push
  - Create cmd/publisher/main.go: poll outbox, push to Redis events list
  - Verify: submit job, see event flow API → outbox → publisher → Redis
  - Verify: kill publisher mid-batch, restart, all events eventually published

Weekend 3 — Consumer
  - Migration 003: create processed_events table
  - Create cmd/consumer/main.go: BLPOP events list, idempotency check, log
  - Verify: events flow end-to-end
  - Verify: republished events (force by killing publisher) are skipped by consumer
  - Document known consumer-crash limitation

Weekend 4 — Reflect and bridge to SD prep
  - Write architectural decision record across V1-V4
  - First mock SD interview using this project as reference
  - Identify gaps to focus SD prep on
```

---

## What you can discuss after V4

**"What's the dual-write problem?"**
Two writes that should be atomic but can't be. The outbox solves it by making both writes go through one DB transaction, then asynchronously bridging to the queue.

**"What is the transactional outbox pattern?"**
Write the event to an outbox table in the same transaction as the state change. A separate publisher reads the outbox and forwards events. Guarantees the event is published if and only if the state change is committed.

**"What is at-least-once delivery and what does it require?"**
The publisher may republish events after a crash. Consumers must be idempotent — able to receive the same event multiple times without changing the outcome. Achieved via deduplication tables, idempotency keys, or transactional consumer logic.

**"Why split services?"**
Independent scaling, independent deployment, fault isolation, clear boundaries. Costs: more operational surface, network failures between services, distributed debugging.

**"What are the costs of a shared database between services?"**
Schema coupling — one service's migration affects others. Contention on hot tables. Implicit coordination that can be hard to reason about. Eventually leads to splitting the database per service.

**"How do you handle duplicate events on the consumer side?"**
Idempotency key (event ID), deduplication table, check-before-process. The consumer treats receipt as the message, not the action — only acts on first receipt.

---

## What V4 does NOT do

- Multiple Redis or Postgres instances
- Cross-region replication
- Real meaningful business logic in the consumer
- Schema versioning for events
- Consumer ACK protocol or transactional read-and-process
- Metrics, tracing, distributed observability
- Outbox cleanup (events accumulate; cleanup is mentioned as a separate concern but not implemented)

V4 is the final version. After this, the project is complete and the focus shifts to SD interview preparation.
