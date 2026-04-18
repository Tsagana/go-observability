package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"go-observability/internal/job"
)

type Dispatcher struct {
	store        job.Storer
	jobs         chan *job.Job
	workerCount  int
	pollInterval time.Duration
	wg           sync.WaitGroup
}

func NewDispatcher(store job.Storer, workerCount int, bufferSize int, pollInterval time.Duration) *Dispatcher {
	return &Dispatcher{
		store:        store,
		jobs:         make(chan *job.Job, bufferSize),
		workerCount:  workerCount,
		pollInterval: pollInterval,
	}
}

// Run starts the worker pool and poll loop, blocks until ctx is cancelled.
func (d *Dispatcher) Run(ctx context.Context) {
	for i := 0; i < d.workerCount; i++ {
		d.wg.Add(1)
		go d.runWorker(ctx, i)
	}
	d.poll(ctx)
	close(d.jobs)
	d.wg.Wait()
}

func (d *Dispatcher) poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			job, err := d.store.Claim(ctx)
			if err != nil || job == nil {
				time.Sleep(d.pollInterval)
				continue
			}
			select {
			case d.jobs <- job:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (d *Dispatcher) runWorker(ctx context.Context, id int) {
	defer d.wg.Done()

	for job := range d.jobs {
		log.Printf("worker %d processing job %s", id, job.ID)
		res, err := process(job)
		writeCtx := context.WithoutCancel(ctx)
		if err != nil {
			_, err = d.store.Fail(writeCtx, job)
			if err != nil {
				log.Println("store.Fail error:", err)
			}
			continue
		}
		err = d.store.Complete(writeCtx, job.ID, res)
		if err != nil {
			log.Println("store.Complete error:", err)
		}
	}
}
