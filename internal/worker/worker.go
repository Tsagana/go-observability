package worker

import (
	"context"
	"log"
	"time"

	"go-observability/internal/job"
)

type Worker struct {
	store        job.Storer
	pollInterval time.Duration
	processFn    func(*job.Job) ([]byte, error)
}

func New(store job.Storer, pollInterval time.Duration) *Worker {
	return &Worker{store: store, pollInterval: pollInterval, processFn: process}
}

func (w *Worker) Run(ctx context.Context) {
	log.Println("worker started")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.processNext(ctx)
			time.Sleep(w.pollInterval)
		}
	}
}

func (w *Worker) processNext(ctx context.Context) {
	job, err := w.store.Claim(ctx)
	if err != nil {
		log.Println("store.Claim error:", err)
		return
	}
	if job == nil {
		return
	}

	res, err := w.processFn(job)

	if err != nil {
		_, err = w.store.Fail(ctx, job)
		if err != nil {
			log.Println("store.Fail error:", err)
		}
		return
	}
	err = w.store.Complete(ctx, job.ID, res)
	if err != nil {
		log.Println("store.Complete error:", err)
	}
}

func process(j *job.Job) ([]byte, error) {
	// TODO: actual processing logic
	log.Printf("processing job %s", j.ID)
	return []byte(`{"result": true}`), nil
}
