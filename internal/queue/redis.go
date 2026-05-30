package queue

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client        *redis.Client
	pendingKey    string
	processingKey string
	claimTimeout  time.Duration
}

func NewRedisQueue(client *redis.Client, pendingKey, processingKey string, claimTimeout time.Duration) *RedisQueue {
	return &RedisQueue{
		client:        client,
		pendingKey:    pendingKey,
		processingKey: processingKey,
		claimTimeout:  claimTimeout,
	}
}

func (q *RedisQueue) Push(ctx context.Context, jobID string) error {
	return q.client.LPush(ctx, q.pendingKey, jobID).Err()
}

// Claim blocks up to claimTimeout waiting for a job. BLMOVE atomically pops
// from pending and pushes to processing — no window where the ID is in neither.
// Returns ("", nil) on timeout so the caller can check ctx and loop.
func (q *RedisQueue) Claim(ctx context.Context) (string, error) {
	id, err := q.client.BLMove(ctx, q.pendingKey, q.processingKey, "LEFT", "RIGHT", q.claimTimeout).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

func (q *RedisQueue) Ack(ctx context.Context, jobID string) error {
	return q.client.LRem(ctx, q.processingKey, 0, jobID).Err()
}

func (q *RedisQueue) ListProcessing(ctx context.Context) ([]string, error) {
	return q.client.LRange(ctx, q.processingKey, 0, -1).Result()
}

func (q *RedisQueue) ListPending(ctx context.Context) ([]string, error) {
	return q.client.LRange(ctx, q.pendingKey, 0, -1).Result()
}
