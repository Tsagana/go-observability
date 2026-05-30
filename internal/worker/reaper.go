package worker

import (
	"context"
	"log/slog"
	"time"

	"go-observability/internal/job"
	"go-observability/internal/queue"
)

type Reaper struct {
	store           job.Storer
	queue           queue.Queue
	interval        time.Duration
	stuckJobTimeout time.Duration
}

func NewReaper(store job.Storer, q queue.Queue, interval time.Duration, stuckJobTimeout time.Duration) *Reaper {
	return &Reaper{store: store, queue: q, interval: interval, stuckJobTimeout: stuckJobTimeout}
}

func (r *Reaper) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(r.interval):
			r.recoverStuck(ctx)
			r.requeueOrphaned(ctx)
		}
	}
}

func (r *Reaper) recoverStuck(ctx context.Context) {
	ids, err := r.queue.ListProcessing(ctx)
	if err != nil {
		slog.Error("reaper.recoverStuck failed to get ids", "error", err)
	}

	for _, id := range ids {
		//load job to read UpdatedAt
		job, err := r.store.Get(ctx, id)
		if err != nil {
			slog.Error("reaper.recoverStuck failed to get job", "job_id", id, "error", err)
			continue
		}

		if time.Since(job.UpdatedAt) <= r.stuckJobTimeout {
			continue
		}
		err = r.store.ResetToPending(ctx, id)
		if err != nil {
			slog.Error("reaper.recoverStuck failed reset to pending", "job_id", id, "error", err)
		}
		r.queue.Push(ctx, id)
		r.queue.Ack(ctx, id)
		slog.Info("reaper.stuck.recovered", "job_id", id)
	}
}

func (r *Reaper) requeueOrphaned(ctx context.Context) {
	postgresIds, err := r.store.GetRetryReadyPending(ctx)
	if err != nil {
		slog.Error("reaper.requeueOrphaned failed to get retry ready ids", "error", err)
		return
	}
	processingIds, err := r.queue.ListProcessing(ctx)
	if err != nil {
		slog.Error("reaper.requeueOrphaned failed to get processing ids", "error", err)
		return
	}
	pendingIds, err := r.queue.ListPending(ctx)
	if err != nil {
		slog.Error("reaper.requeueOrphaned failed to get pending ids", "error", err)
		return
	}
	redisIds := make(map[string]bool, len(processingIds)+len(pendingIds))
	for _, id := range processingIds {
		redisIds[id] = true
	}
	for _, id := range pendingIds {
		redisIds[id] = true
	}
	for _, id := range postgresIds {
		if !redisIds[id] {
			r.queue.Push(ctx, id)
			slog.Info("reaper.orphan.requeued", "job_id", id)
		}
	}
}
