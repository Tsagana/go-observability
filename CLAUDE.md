# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the service (Docker Compose)
make up
make down
make logs

# Database migrations
make migrate-up
make migrate-down
make migrate-force VERSION=<n>

# Tests
go test ./...                                        # all tests
go test ./internal/job/                              # single package
go test -run TestCalculateBackoff ./internal/job/   # single test
go test -v ./...                                     # verbose
```

No linter is configured. Build happens inside the Docker container (`go build -o server ./cmd/server`).

## Architecture

This is an **async job processing service** — clients submit jobs via HTTP, workers claim and process them from PostgreSQL, and clients poll for results.

**Data flow:**
```
POST /jobs → API handler → job_store.Create() → PostgreSQL (status: pending)
                                ↓
                    Dispatcher polls DB every POLL_INTERVAL
                    Claims job atomically (SELECT FOR UPDATE SKIP LOCKED)
                    Sends job to worker pool via buffered channel
                                ↓
                    Worker executes (executor.go, 15-sec simulated task)
                    On success: updates job to completed + writes result JSONB
                    On failure: retries with exponential backoff (max 3 retries)
                                ↓
GET /jobs/{id} → returns current job state
```

**Package layout:**
- `cmd/server/` — entry point; wires up DB, store, dispatcher, reaper, HTTP server
- `internal/api/` — HTTP handlers and routes for `POST /jobs`, `GET /jobs/{id}`, `GET /healthz`
- `internal/job/` — `Job` domain model, status constants, `Store` interface, PostgreSQL queries
- `internal/worker/` — `Dispatcher` (manages worker pool + polling loop), `Reaper` (recovers stuck jobs), `executor.go` (job processing logic)
- `internal/db/` — PostgreSQL connection pool setup

**Key design decisions:**
- `SELECT FOR UPDATE SKIP LOCKED` ensures at-most-once job claiming under concurrent workers
- No message queue — the DB is the queue (intentional V1 simplicity)
- Reaper periodically resets jobs stuck in `processing` state (configurable timeout)
- Retry backoff: `2^retry_count` seconds, capped at `MaxRetry = 3`

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | (required) | PostgreSQL connection string |
| `WORKER_COUNT` | 5 | Concurrent workers |
| `JOB_CHANNEL_BUFFER` | = WORKER_COUNT | Buffered channel size |
| `POLL_INTERVAL` | 2s | DB polling frequency |
| `STUCK_JOB_TIMEOUT` | 300s | Age before a processing job is considered stuck |
| `REAPER_INTERVAL` | 60s | How often the reaper runs |
| `PORT` | 8080 | HTTP listen port |

## Database

Single `jobs` table with UUID primary key, `status` (`pending`/`processing`/`completed`/`failed`), `payload` JSONB input, `result` JSONB output, `error` text, `retry_count`, `retry_after`, and standard timestamps. Index on `(status, created_at)` for the dispatcher's polling query.

Migrations live in `migrations/` and are managed by `golang-migrate`.
