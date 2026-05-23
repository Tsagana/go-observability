# Architecture: Job Processing System — V3

## What V3 is

Replace the dispatcher's polling loop with a Redis-backed queue. Postgres stays the source of truth for state. Redis becomes the delivery mechanism.

One focused change. The rest of the system is unchanged.

---

## What does not change

- Postgres remains source of truth for state
- `jobs` table schema unchanged
- State machine unchanged: pending → processing → completed / failed
- API endpoints unchanged
- Worker pool, channel, executor, AI loop unchanged
- Retry logic, backoff, error classification unchanged

---

## What changes

### POST /jobs
After writing to Postgres, push job ID to Redis pending list.

```
1. Store.Create() → INSERT, status=pending
2. queue.Push(jobID) → LPUSH to Redis pending list
3. Return 202
```

If step 2 fails after step 1 succeeds, the reaper recovers it later.

### Dispatcher
No more SQL polling. Block on Redis.

```
1. queue.Claim() → BLMOVE pending → processing (atomic, blocking)
2. Store.Get(jobID) → load full job from Postgres
3. Store.MarkProcessing(jobID)
4. Send to worker channel
```

BLMOVE is atomic pop-and-push. No race condition, no SKIP LOCKED needed in hot path.

### Worker completion
After writing result to Postgres, ack the job in Redis.

```
On any terminal outcome (success, retry, permanent fail):
  - Update Postgres (Store.Complete / Fail / FailPermanently)
  - queue.Ack(jobID) → LREM from Redis processing list
```

### Reaper
Three recovery scenarios now:

1. **Stuck in processing** — in Redis processing list + Postgres `updated_at` stale → reset to pending, requeue
2. **Retry-after expired** — Postgres pending + `retry_after < NOW()` → push to Redis
3. **Orphaned by failed push** — Postgres pending, not in Redis → push to Redis

### Shutdown
BLMOVE must be cancellable via context. Use finite timeout (5s) so dispatcher periodically checks for cancellation.

---

## Redis structures

```
jobs:pending      list of job IDs waiting to be claimed
jobs:processing   list of job IDs currently claimed
```

Job IDs only. Payloads live in Postgres.

---

## Queue interface

```go
type Queue interface {
    Push(ctx, jobID) error
    Claim(ctx) (jobID, error)         // blocking
    Ack(ctx, jobID) error
    ListProcessing(ctx) ([]string, error)
}
```

Two implementations: `RedisQueue` (production) and `MemoryQueue` (tests).

---

## New config

```
REDIS_URL=redis://localhost:6379
QUEUE_PENDING_KEY=jobs:pending
QUEUE_PROCESSING_KEY=jobs:processing
QUEUE_CLAIM_TIMEOUT=5s
```

---

## Project structure additions

```
internal/
  queue/
    queue.go        Queue interface
    redis.go        RedisQueue
    memory.go       MemoryQueue (tests)
  redis/
    client.go       Redis connection
```

---

## Key decisions

**Why hybrid Postgres + Redis?**
Postgres for durability and state. Redis for fast push-based delivery. Standard pattern (Sidekiq, Resque).

**Why BLMOVE not BRPOP + LPUSH?**
Atomic. No window where the job ID is in neither list.

**Why job IDs only?**
Bounded Redis memory. No serialization mismatch. Postgres stays unambiguous source of truth.

**Why finite BLMOVE timeout?**
Lets dispatcher check context cancellation. 5s is acceptable shutdown latency.

**Why not solve dual-write atomically here?**
Reaper covers the gap. Atomic dual-write is V4 (outbox pattern).

---

## Failure modes

| Failure | Recovery |
|---|---|
| API write OK, Redis push fails | Reaper sweep finds orphaned pending jobs |
| Worker crashes mid-job | Reaper detects stuck via Postgres `updated_at` |
| Worker completes, Ack fails | Reaper sees terminal status in Postgres, removes from Redis |
| Redis crashes, loses data | Reaper rebuilds from Postgres state |
| Postgres unavailable | API returns 5xx, workers stop, recover when DB returns |

---

## Build order

```
1. Add Redis to docker-compose, create redis client
2. Queue interface + Redis + Memory implementations + tests
3. Wire POST /jobs to push after Store.Create
4. Replace dispatcher poll loop with BLMOVE
5. Wire worker completion to call Ack
6. Update reaper with 3 recovery scenarios
7. Verify graceful shutdown with in-flight jobs
8. Load test: compare latency and DB load vs V2
```
