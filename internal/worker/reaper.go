package worker

import (
	"context"
	"log/slog"
	"time"

	"go-observability/internal/job"
)

type Reaper struct {
	store           job.Storer
	interval        time.Duration
	stuckJobTimeout time.Duration
}

func NewReaper(store job.Storer, interval time.Duration, stuckJobTimeout time.Duration) *Reaper {
	return &Reaper{store: store, interval: interval, stuckJobTimeout: stuckJobTimeout}
}

func (r *Reaper) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(r.interval):
			count, err := r.store.RecoverStuck(ctx, r.stuckJobTimeout)
			if err != nil {
				slog.Error("store.recoverStuck failed", "error", err)
			}
			slog.Info("reaper.recovered", "count", count)
		}
	}
}
