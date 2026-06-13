package events

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
)

// Queue is the read side of the events Redis list.
type Queue interface {
	Pop(ctx context.Context) (string, error)
}

// EventEnvelope is the shape pushed onto the events Redis list by the publisher.
type EventEnvelope struct {
	EventID   int64  `json:"event_id"`
	EventType string `json:"event_type"`
	JobID     string `json:"job_id"`
}

type Consumer struct {
	store *Store
	queue Queue
}

func NewConsumer(store *Store, queue Queue) *Consumer {
	return &Consumer{
		store: store,
		queue: queue,
	}
}

// Start runs the consume loop until ctx is cancelled.
func (c *Consumer) Start(ctx context.Context) error {
	for {
		result, err := c.queue.Pop(ctx)
		if err != nil {
			if isContextError(err) {
				return err
			}
			slog.Error("consumer.start failed", "error", err)
			continue
		}
		if result == "" {
			continue
		}

		envelope, err := decodeEnvelope(result)
		if err != nil {
			slog.Error("consumer.decodeEnvelope failed", "error", err, "data", result)
			continue
		}

		marked, err := c.store.MarkProcessed(ctx, envelope.EventID)
		if err != nil {
			slog.Error("consumer.MarkProcessed failed", "error", err, "event_id", envelope.EventID)
			continue
		}
		if marked {
			slog.Info("consumer.MarkProcessed - event processed", "event_id", envelope.EventID, "event_type", envelope.EventType, "job_id", envelope.JobID)
		} else {
			slog.Info("consumer.MarkProcessed - duplicate skipping", "event_id", envelope.EventID, "event_type", envelope.EventType, "job_id", envelope.JobID)

		}
	}
}

func decodeEnvelope(data string) (EventEnvelope, error) {
	var e EventEnvelope
	return e, json.Unmarshal([]byte(data), &e)
}

func isContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

// keep the compiler happy until Start is implemented
var _ = slog.Info
