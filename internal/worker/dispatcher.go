package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"go-observability/internal/ai"
	"go-observability/internal/job"
	"go-observability/internal/queue"
)

type Dispatcher struct {
	store       job.Storer
	jobs        chan *job.Job
	workerCount int
	queue       queue.Queue
	wg          sync.WaitGroup
	aiClient    *ai.Client
}

func NewDispatcher(store job.Storer, workerCount int, bufferSize int, queue queue.Queue, aiClient *ai.Client) *Dispatcher {
	return &Dispatcher{
		store:       store,
		jobs:        make(chan *job.Job, bufferSize),
		workerCount: workerCount,
		queue:       queue,
		aiClient:    aiClient,
	}
}

// Run starts the worker pool and poll loop, blocks until ctx is cancelled.
func (d *Dispatcher) Run(ctx context.Context) {
	for i := 0; i < d.workerCount; i++ {
		d.wg.Add(1)
		go d.runWorker(ctx, i)
		slog.Info("worker.started", "worker_id", i)
	}
	d.claim(ctx)
	close(d.jobs)
	d.wg.Wait()
}

func (d *Dispatcher) claim(ctx context.Context) {
	for {
		id, err := d.queue.Claim(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			slog.Error("dispatcher.claim failed", "error", err)
			continue
		}
		if id == "" {
			continue
		}

		job, err := d.store.Get(ctx, id)
		if err != nil {
			slog.Error("store.get failed", "error", err)
			continue
		}
		err = d.store.MarkProcessing(ctx, id)
		if err != nil {
			slog.Error("store.markProcessing failed", "error", err)
			continue
		}
		select {
		case d.jobs <- job:
		case <-ctx.Done():
			return
		}
	}
}

func (d *Dispatcher) runWorker(ctx context.Context, id int) {
	defer d.wg.Done()
	defer slog.Info("worker.stopped", "worker_id", id)

	for j := range d.jobs {
		slog.Info("job.claimed", "job_id", j.ID, "worker_id", id)
		timeBeforeStart := time.Now()
		res, err := process(ctx, j, d.aiClient)

		writeCtx := context.WithoutCancel(ctx)
		if err != nil {
			if job.IsRetryable(err) {
				slog.Error("job.failed", "job_id", j.ID, "worker_id", id, "error", err)
				_, err = d.store.Fail(writeCtx, j)
				if err != nil {
					slog.Error("store.fail failed", "error", err)
				}
				err = d.queue.Ack(writeCtx, j.ID)
				if err != nil {
					slog.Error("queue.ack failed", "error", err)
				}
				continue
			} else {
				slog.Error("job.failed.permanent", "job_id", j.ID, "worker_id", id, "error", err)
				err = d.store.FailPermanently(writeCtx, j.ID, err.Error())
				if err != nil {
					slog.Error("store.fail failed", "error", err)
				}
				err = d.queue.Ack(writeCtx, j.ID)
				if err != nil {
					slog.Error("queue.ack failed", "error", err)
				}
				continue
			}
		}
		err = d.store.Complete(writeCtx, j.ID, res)

		if err != nil {
			slog.Error("store.complete failed", "error", err)
		} else {
			duration := time.Since(timeBeforeStart).Milliseconds()
			slog.Info("job.completed", "job_id", j.ID, "worker_id", id, "duration_ms", duration)
		}
		err = d.queue.Ack(writeCtx, j.ID)
		if err != nil {
			slog.Error("queue.ack failed", "error", err)
		}
	}
}
