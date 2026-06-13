package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"go-observability/internal/events"
	"go-observability/internal/queue"
)

type Publisher struct {
	store        *Store
	jobsQueue    queue.Queue
	eventsQueue  queue.Queue
	batchSize    int
	pollInterval time.Duration
}

func NewPublisher(store *Store, jobsQueue queue.Queue, eventsQueue queue.Queue, batchSize int, pollInterval time.Duration) *Publisher {
	return &Publisher{
		store:        store,
		jobsQueue:    jobsQueue,
		eventsQueue:  eventsQueue,
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

// publishBatch fetches one batch of unpublished events inside a single transaction,
// pushes each event to both queues, then marks them published before committing.
func (p *Publisher) publishBatch(ctx context.Context) (int, error) {
	tx, err := p.store.BeginTx(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	batch, err := p.store.FetchUnpublishedTx(ctx, tx, p.batchSize)
	if err != nil {
		return 0, err
	}

	var ids []int64
	for _, event := range batch {
		payload, err := parsePayload(event.Payload)
		if err != nil {
			slog.Error("publisher.parsePayload failed", "event_id", event.ID, "error", err)
			continue
		}

		if err := p.jobsQueue.Push(ctx, payload.JobID); err != nil {
			slog.Error("publisher.jobsQueue.Push failed", "event_id", event.ID, "error", err)
			continue
		}

		envelope, err := json.Marshal(events.EventEnvelope{
			EventID:   event.ID,
			EventType: event.EventType,
			JobID:     payload.JobID,
		})
		if err != nil {
			slog.Error("publisher.marshal envelope failed", "event_id", event.ID, "error", err)
			continue
		}

		if err := p.eventsQueue.Push(ctx, string(envelope)); err != nil {
			slog.Error("publisher.eventsQueue.Push failed", "event_id", event.ID, "error", err)
			continue
		}

		ids = append(ids, event.ID)
	}

	if len(ids) == 0 {
		return 0, nil
	}

	if err := p.store.MarkPublishedTx(ctx, tx, ids); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return len(ids), nil
}

// jobCreatedPayload is the shape written into outbox.payload by the API handler.
type jobCreatedPayload struct {
	JobID string `json:"job_id"`
}

func parsePayload(raw []byte) (jobCreatedPayload, error) {
	var p jobCreatedPayload
	return p, json.Unmarshal(raw, &p)
}
