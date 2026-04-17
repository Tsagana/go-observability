package worker

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"go-observability/internal/job"
)

type mockStore struct {
	claimFn    func(ctx context.Context) (*job.Job, error)
	completeFn func(ctx context.Context, id string, result []byte) error
	failFn     func(ctx context.Context, j *job.Job) (*job.Job, error)
}

func (m *mockStore) Claim(ctx context.Context) (*job.Job, error) {
	return m.claimFn(ctx)
}

func (m *mockStore) Complete(ctx context.Context, id string, result []byte) error {
	return m.completeFn(ctx, id, result)
}

func (m *mockStore) Fail(ctx context.Context, j *job.Job) (*job.Job, error) {
	return m.failFn(ctx, j)
}

func TestProcessNext_Success(t *testing.T) {
	completeCalled := false
	j := &job.Job{ID: "test-id", Status: job.StatusPending}

	store := &mockStore{
		claimFn: func(ctx context.Context) (*job.Job, error) {
			return j, nil
		},
		completeFn: func(ctx context.Context, id string, result []byte) error {
			completeCalled = true
			return nil
		},
		failFn: func(ctx context.Context, j *job.Job) (*job.Job, error) {
			t.Error("Fail should not be called on success")
			return nil, nil
		},
	}

	w := newTestWorker(store)
	w.processNext(context.Background())

	if !completeCalled {
		t.Error("Complete was not called")
	}
}

func TestProcessNext_ProcessFails_FailCalled(t *testing.T) {
	failCalled := false
	j := &job.Job{ID: "test-id", Status: job.StatusPending}

	store := &mockStore{
		claimFn: func(ctx context.Context) (*job.Job, error) {
			return j, nil
		},
		completeFn: func(ctx context.Context, id string, result []byte) error {
			t.Error("Complete should not be called on failure")
			return nil
		},
		failFn: func(ctx context.Context, j *job.Job) (*job.Job, error) {
			log.Println("Job failed ", j.ID)
			failCalled = true
			return j, nil
		},
	}

	w := newTestWorker(store)
	w.processFn = func(j *job.Job) ([]byte, error) {
		return nil, errors.New("simulated failure")
	}
	w.processNext(context.Background())

	if !failCalled {
		t.Error("Fail was not called")
	}
}

func newTestWorker(store *mockStore) *Worker {
	return New(store, time.Millisecond)
}
