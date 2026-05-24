package queue

import "context"

type Queue interface {
	// Push adds a job ID to the pending list.
	Push(ctx context.Context, jobID string) error

	// Claim blocks until a job is available, then atomically moves it from
	// pending to processing. Returns ("", nil) on timeout — not an error.
	Claim(ctx context.Context) (string, error)

	// Ack removes a job ID from the processing list.
	Ack(ctx context.Context, jobID string) error

	// ListProcessing returns all job IDs currently in the processing list.
	ListProcessing(ctx context.Context) ([]string, error)
}
