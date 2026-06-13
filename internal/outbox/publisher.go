package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"go-observability/internal/queue"
)

type Publisher struct {
	store        *Store
	jobsQueue    queue.Queue
	batchSize    int
	pollInterval time.Duration
}

func NewPublisher(store *Store, jobsQueue queue.Queue, batchSize int, pollInterval time.Duration) *Publisher {
	return &Publisher{
		store:        store,
		jobsQueue:    jobsQueue,
		batchSize:    batchSize,
		pollInterval: pollInterval,
	}
}

// Start runs the publish loop until ctx is cancelled.
func (p *Publisher) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(p.pollInterval):
			_, err := p.publishBatch(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return err
				}
				slog.Error("publisher.publishBatch failed", "error", err)
				continue
			}
		}
	}
}

// publishBatch fetches one batch of unpublished events, pushes job IDs to Redis,
// and marks them published. Returns the number of events published.
func (p *Publisher) publishBatch(ctx context.Context) (int, error) {
	events, err := p.store.FetchUnpublished(ctx, p.batchSize)
	if err != nil {
		slog.Error("store.FetchUnpublished failed", "error", err)
		return 0, err
	}
	var ids []int64
	for _, event := range events {
		jobID, err := parseJobID(event.Payload)
		if err != nil {
			slog.Error("publisher.parseJobID failed", "event_id", event.ID, "error", err)
			continue
		}
		err = p.jobsQueue.Push(ctx, jobID)
		if err != nil {
			slog.Error("jobsQueue.push failed", "event_id", event.ID, "error", err)
			continue
		}
		ids = append(ids, event.ID)
	}
	err = p.store.MarkPublished(ctx, ids)
	if err != nil {
		slog.Error("store.markPublished push failed", "error", err)
		return 0, err
	}
	return len(ids), nil
}

// jobCreatedPayload is the shape written into outbox.payload by the API handler.
type jobCreatedPayload struct {
	JobID string `json:"job_id"`
}

func parseJobID(payload []byte) (string, error) {
	var p jobCreatedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return "", err
	}
	return p.JobID, nil
}

