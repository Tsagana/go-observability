Architectural decision record (ADR) 14:52

# Polling vs queue (V2 → V3)
My answer:
Problem: initially we had a pool of workers that would hit db to get next job. It was decided to move workers to get data on jobs in queue instead and leave db only as source of truth.
Choice: use Redis queue instead of db to claim jobs
Tradeoff: support and maintain Redis as another entity can make infra more complex and more difficult to maintain. Redis can store incorrect data if db went out of service(?)

AI:
Problem: how workers discover new jobs to process without a dedicated message broker. 
Choice: keep PostgreSQL as the queue — the dispatcher polls the jobs table every POLL_INTERVAL and claims work atomically via SELECT FOR UPDATE SKIP LOCKED, rather than introducing Kafka/RabbitMQ/SQS with push-based delivery. 
The tradeoff: this keeps the system simple (one datastore, no extra infra to deploy/monitor, transactional consistency between job state and claiming) but at the cost of latency (jobs sit up to POLL_INTERVAL before being picked up) and eventual scalability limits (polling load and row-lock contention grow with worker count and table size in a way a purpose-built queue wouldn't).

# Two transactions in the worker (claim, then write result)
My answer:
Problem: To ensure that result will be written in the same flow as claiming
Choice: SELECT + UPDATE row in db by same worker
Tradeoff: Worker can crash after claiming and not finish processing

AI:
The problem is that claiming a job (marking it processing) and recording its outcome (marking it completed/failed with result or error) happen at different points in time separated by the actual job execution — potentially seconds of work — so they can't share a single transaction. 
The choice was to use two separate, short-lived transactions: one to atomically claim the job via SELECT FOR UPDATE SKIP LOCKED, and a second, independent one after execution to persist the result. 
The tradeoff: this avoids holding a row lock (and a DB connection) for the entire duration of job execution, which would badly limit concurrency and risk long-running transactions — but it opens a window where a worker can crash after claiming but before writing a result, leaving the job stuck in processing with no completion record, which is exactly why the reaper exists as a compensating mechanism to reclaim orphaned jobs.

# SELECT FOR UPDATE SKIP LOCKED for claiming jobs
My answer:
Problem: To ensure that same job won't be claimed or processed by different worker processes and thus be duplicated
Choice: SELECT FOR UPDATE SKIP LOCKED, worker doesn't claim job already locked
Tradeoff: if mid crashed, it stays locked and can't be claimed anymore(X)

AI:
The problem is how multiple concurrent dispatcher/worker processes can pull pending jobs from the same jobs table without two workers grabbing the same job (a race condition) or blocking each other while waiting on locked rows. 
The choice was Postgres's SELECT ... FOR UPDATE SKIP LOCKED, which lets each polling transaction lock and claim a row while silently skipping any rows already locked by a concurrent claim, giving at-most-once delivery without an external coordination service. 
The tradeoff: this is simple, requires no extra infrastructure, and scales fine at moderate concurrency, but it's a database-specific feature (ties the design to Postgres) and doesn't eliminate contention entirely — under high worker counts or large pending backlogs, many transactions scanning/locking the same index range can still cause throughput to degrade compared to a purpose-built queue with push-based work distribution.

# retry_count + retry_after for backoff
My answer:
Problem: after job fail we don't need to immediately send requests again and again with no limits because it could increase load in case response gets delayed for independent reasons.
Choice: Limit to retry_count and increase retry_after exponentially not to increase load and give time to recover
Tradeoff: still can be failed, or could be delay in case of server is recovered faster.

AI:
The problem is that failed jobs shouldn't be retried immediately (which would hammer a struggling downstream dependency and waste worker cycles) nor retried forever (which would let permanently-broken jobs consume resources indefinitely). 
The choice was to store retry_count and retry_after directly as columns on the jobs row: each failure increments retry_count, computes the next eligible time as 2^retry_count seconds out and stores it in retry_after, and the dispatcher's polling query simply filters on retry_after <= now() alongside status = 'pending', capping retries at MaxRetry = 3 before marking the job permanently failed. 
The tradeoff: this keeps backoff state colocated with the job and queryable via the same simple polling mechanism (no separate scheduler or delay queue needed), but it's coarse — backoff timing is only as precise as POLL_INTERVAL, and retry policy is hardcoded per-job rather than configurable per job type or dynamically tunable without a schema change.

# Stuck job reaper via updated_at timeout
My answer:
Problem: failed jobs can stuck in processing status indefinitely, need to picke them up and process again
Choice: Introduce separate job, which goes into db and check jobs in processing state which updated_at time exceeds some limit that we configured.
Tradeoff: dependent on db, limit is hardcoded

AI:
The problem is that a worker can die mid-job (crash, OOM kill, pod eviction) after claiming a job but before writing its result, leaving that job permanently stuck in processing with no owner left to finish or fail it. 
The choice was a reaper goroutine that runs every REAPER_INTERVAL and resets any job still in processing whose updated_at is older than STUCK_JOB_TIMEOUT back to pending (or failed, depending on retry count), using elapsed wall-clock time as a proxy for "the worker that owned this is gone." 
The tradeoff: this is a simple, infrastructure-free way to recover from crashes without needing worker heartbeats, leases, or a liveness registry — but it's a blunt instrument: STUCK_JOB_TIMEOUT has to be set well above the longest legitimate job duration, which means genuinely stuck jobs sit unrecovered for a while, and a job that's merely slow (not dead) risks being reaped and re-run concurrently with its still-running original if the timeout is set too aggressively.

# Dispatcher + channel + worker pool (vs N independent pollers)
My answer:
Problem: If there are only workers directly hitting db, there could be a thundering herd problem, when same worker tries to claim same jobs at the same time, increasing load. Need to distribute jobs equally.
Choice: Introduce dispatcher, which job is to claim jobs from db and put it to dedicated channel, from which workers would claim their jobs.
Tradeoff: additional infra to support that

AI:
The problem is how to bound and coordinate concurrency — with WORKER_COUNT workers needing to execute jobs, you could have each worker independently poll the DB, but that means N separate polling queries hitting Postgres on every interval, with no natural backpressure between "jobs being claimed" and "jobs actually able to run." 
The choice was a single dispatcher that owns the polling loop, claims jobs, and pushes them onto a buffered Go channel (JOB_CHANNEL_BUFFER), with a fixed pool of worker goroutines pulling from that channel to execute them. 
The tradeoff: this centralizes and simplifies the polling logic (one query pattern, one place to reason about claim correctness) and gives free backpressure via the channel buffer (the dispatcher naturally slows down claiming once workers are saturated), but it introduces a single dispatcher as a serialization point/potential bottleneck and a bit more in-process coordination complexity (channel sizing, dispatcher lifecycle, graceful shutdown of both dispatcher and pool) compared to N fully independent, symmetric pollers.

# Why hybrid Postgres + Redis (vs Redis-only)
My answer:
Problem: Need to have one source of truth which will record jobs claimed, processing and in such way they can be re-processed some failed jobs.
Wastes cycles.
Choice: Use PG as source of truth, Redis as a source of processing jobs
Tradeoff: Introduces dual write problem

AI:
The problem is that V2's DB-polling dispatcher (SELECT FOR UPDATE SKIP LOCKED on an interval) wastes cycles polling an empty table and adds up to POLL_INTERVAL of latency between a job being created and being claimed — but going all-in on Redis would mean losing durable, transactionally-consistent state. 
The choice was to split responsibilities: Postgres stays the single source of truth for job state (still owns status, retry_count, results, and survives restarts/crashes with full durability), while Redis becomes a thin, fast delivery layer — LPUSH/BLMOVE on job IDs only, no job payloads or results stored there — so the dispatcher blocks on Redis instead of polling Postgres. 
The tradeoff: this gets push-based, low-latency delivery without sacrificing durability (Redis can lose its data entirely and the reaper rebuilds queue state from Postgres), but it introduces a dual-write problem — the API now writes to two systems (INSERT then LPUSH) that aren't atomic, creating a window where a job lands in Postgres but never makes it into the Redis queue, recoverable only via the reaper's orphan-sweep (a gap V4's outbox pattern later closes).

# Job IDs only in Redis, payloads in Postgres
My answer:
Problem: To keep state of jobs same for everyone, keep PG as source of truth and redis lightweight
Choice: store only job ids in redis and retrieve data/write by it
Tradeoff: More load on db

AI:
The problem is what to actually put in the Redis queue structures (jobs:pending, jobs:processing) — pushing full job payloads into Redis would mean two copies of job data to keep in sync, plus serialization concerns (schema evolution, payload size limits, JSONB fields that don't map cleanly onto Redis's simpler value types). 
The choice was to keep Redis holding only job IDs (UUIDs) and treat it purely as a delivery/ordering mechanism, with the dispatcher doing a Store.Get(jobID) back to Postgres to load the full job after claiming it off the queue. 
The tradeoff: this keeps Redis memory usage small and bounded regardless of payload size, avoids any serialization mismatch between two representations of the same job, and keeps Postgres unambiguously the single source of truth — at the cost of an extra DB round-trip per claimed job (fetching the full row after the queue pop) that a payload-in-Redis design would have avoided.

# BLMOVE over BRPOP+LPUSH
My answer:
Problem: Make it in one transaction or split into two, while processing jobs.
Choice: use BLMOVE only so it does pop and push inside single transaction without failing in between
Tradeoff: idk, Redis specific

AI:
The problem is how the dispatcher atomically moves a job ID from the jobs:pending list to jobs:processing when claiming it — doing this as two separate Redis commands (BRPOP to pop from pending, then LPUSH to push onto processing) leaves a window between the two calls where, if the dispatcher crashes, the job ID exists in neither list and is silently lost. 
The choice was BLMOVE, Redis's atomic blocking pop-and-push primitive, which performs both halves of the move as a single indivisible operation. 
The tradeoff: this closes the crash window entirely with no added complexity in the queue code (it's a single Redis call, not a transaction or Lua script), but it does tie the implementation to a Redis feature specific to this pop-then-push pattern — a general two-queue move isn't expressible this cleanly on every backend, so swapping the delivery layer later would mean re-deriving this atomicity guarantee some other way.

# Service split (API / worker / publisher / consumer)
My answer:
Problem: To move from monolit to microservices architecture, so each service could be scaled independently
Choice: create different services
Tradeoff: Maintenance, additional infra, fail in betweens, needs consistency mechanism

AI:
The problem is that up through V3, the API and worker ran in a single process — sharing a DB and Redis but no in-process state — which meant they could only scale and deploy together, and any coordination temptation (shared memory, direct function calls) risked coupling that undermines the "coordinate only through infrastructure" discipline. 
The choice was to split into four independent binaries/containers — api (HTTP only), worker (dispatcher + pool + reaper), publisher (outbox → Redis), and consumer (Redis → processed events) — each owning one responsibility and talking to the others only via Postgres and Redis, never directly. 
The tradeoff: this buys independent scaling (e.g., scale worker replicas without touching the API) and independent deployability, and forces the system to prove its coordination actually works through real infrastructure rather than implicit in-process guarantees — but it multiplies operational surface area (four things to deploy, monitor, and version instead of one) and turns previously-atomic in-process sequences into multi-process flows that need their own consistency mechanisms (which is exactly why the outbox pattern shows up alongside this split).

# Transactional outbox for atomic dual-write
My answer:
Problem: Dual write problem if one job are claimed by different workers - wrong
Choice: use outbox pattern
Tradeoff: support consistency and introduces complexity

AI:
The problem is the one V3 explicitly punted on — POST /jobs writes to two systems (INSERT into Postgres, LPUSH to Redis) that aren't atomic, so a Redis push failure after a successful Postgres insert orphans the job until the reaper's sweep eventually catches it, which works but leaves a real window of inconsistency and an extra recovery path to reason about. 
The choice was the transactional outbox pattern: the API writes an outbox row in the same Postgres transaction as the job insert (so both succeed or neither does, by ordinary ACID guarantees), and a separate publisher process polls the outbox table and forwards unpublished events to Redis, marking them published after each successful push. 
The tradeoff: this eliminates the inconsistency window entirely and gets a free audit log of every event as a side effect — but it trades a same-request synchronous push for asynchronous, at-least-once delivery (the publisher can crash between pushing to Redis and marking the row published, so it may republish, meaning the consumer must dedupe), and it adds a whole new process plus table to operate.
# Idempotent consumer via processed_events table
My answer:
Problem: How consumer must take job data from
Choice: add separate processed events, which will have data on consumer job events
Tradeoff: adds complexity and table

AI:
The problem is that the outbox publisher gives only at-least-once delivery — it can crash after pushing an event to Redis but before marking it published, causing the same event to be pushed again on the next iteration — so the consumer on the other end can receive duplicates and needs a way to process each event's effect exactly once despite that. 
The choice was a processed_events table keyed by event_id (a bare primary key on the ID, no payload) that the consumer checks before acting on an event and inserts into right after: if the ID is already present, it logs and skips; the primary-key constraint itself is what enforces idempotency. 
The tradeoff: this is a minimal, cheap way to dedupe against at-least-once delivery using nothing but a Postgres uniqueness constraint — but as implemented it's not fully safe for consumers with real side effects, since "process the event" and "insert into processed_events" aren't in the same transaction, so a crash between the two still allows a duplicate side effect (acceptable here only because "processing" is just a log line; a real consumer would need to wrap both in one transaction).

# Retryable vs permanent error classification in AI loop
My answer:
Problem: Some errors can be retryable and must be processed again vs some jobs are permanently failed and shouldnt be touched anymore
Choice: Classify errors into retryable/permanent and set according statuses. Make reaper re-process retryable jobs.
Tradeoff: ??

AI:
The problem is that the agent loop (internal/ai/client.go) can fail for very different reasons — an invalid tool schema or exceeding max_steps is a configuration/logic bug that will never succeed no matter how many times you retry, while an Anthropic API call failing on a transient rate limit or network blip might succeed on the next attempt — and treating all failures the same way either wastes retries on unfixable errors or gives up too early on recoverable ones. 
The choice was two explicit wrapper error types, RetryableError and PermanentError (job.NewRetryable / job.NewPermanent, internal/job/errors.go), with classification done at the point an error is raised: schema validation and max-steps-exceeded are wrapped as Permanent immediately, while API-call failures are checked via job.IsRetryable and, if so, retried in-loop with exponential backoff up to maxRetries before surfacing to the dispatcher, which checks the same classification to decide between DB-level retry (retry_count/retry_after) and FailPermanently. 
The tradeoff: this lets each error carry its own retry semantics all the way from the AI client through the dispatcher instead of one blanket policy — but IsRetryable defaults to false for any error not explicitly wrapped as RetryableError, so classification correctness lives entirely at each call site's discretion, and an error type someone forgets to classify as retryable (or a raw, unwrapped Anthropic SDK error) silently becomes permanent rather than getting the benefit of the doubt.

# Job-level + step-level timeouts in agent loop
My answer:
Problem: Retry API calls in case if server was not accessible, having problems on its side. In agentics calls, tools needs to be used after making a call and fail in each step
Choice: Introduce timeouts on job level and step level in agentic loop
Tradeoff: Delay in results delivery

AI:
The problem is that an agentic loop calling an LLM in a tool-use cycle (internal/ai/client.go) has two different ways to run too long: a single step (one API call) can hang or take unexpectedly long, and separately, the loop as a whole can burn through many legitimate-but-slow steps without any individual step being the culprit — a single timeout covering the whole loop can't distinguish "one stuck call" from "many steps, each fine, but overall too slow," and using only a per-step timeout can't bound total job duration at all. 
The choice was two nested context.WithTimeout calls: jobTimeout wraps the entire RunAgentLoop invocation once, and stepTimeout wraps each individual step inside the loop (recreated fresh every iteration via stepCtx, stepCancel := context.WithTimeout(ctx, c.stepTimeout)), so a step is bounded independently while the parent job context still caps the overall run. 
The tradeoff: this gives precise failure diagnosis (a step that hangs times out fast without waiting for the full job budget, while a job that's making steady progress but has many steps still gets cut off at a sane ceiling) — but it adds two timeout values to tune instead of one, and getting them wrong in either direction has different failure modes: too-tight stepTimeout kills slow-but-valid API calls, while too-loose jobTimeout lets a multi-step loop run far longer than a worker slot should reasonably be held.