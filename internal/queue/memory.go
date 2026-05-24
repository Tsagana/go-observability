package queue

import (
	"context"
	"sync"
)

type MemoryQueue struct {
	pending    chan string
	mu         sync.Mutex
	processing []string
}

func NewMemoryQueue(bufferSize int) *MemoryQueue {
	return &MemoryQueue{
		pending:    make(chan string, bufferSize),
		processing: []string{},
	}
}

func (q *MemoryQueue) Push(ctx context.Context, jobID string) error {
	select {
	case q.pending <- jobID:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *MemoryQueue) Claim(ctx context.Context) (string, error) {
	select {
	case id := <-q.pending:
		q.mu.Lock()
		q.processing = append(q.processing, id)
		q.mu.Unlock()
		return id, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (q *MemoryQueue) Ack(ctx context.Context, jobID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, id := range q.processing {
		if id == jobID {
			q.processing = append(q.processing[:i], q.processing[i+1:]...)
			return nil
		}
	}
	return nil
}

func (q *MemoryQueue) ListProcessing(ctx context.Context) ([]string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]string, len(q.processing))
	copy(out, q.processing)
	return out, nil
}
