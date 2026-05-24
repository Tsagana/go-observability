package queue_test

import (
	"context"
	"testing"
	"time"

	"go-observability/internal/queue"
)

func TestMemoryQueue_PushClaim(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	ctx := context.Background()

	if err := q.Push(ctx, "job-1"); err != nil {
		t.Fatal(err)
	}

	id, err := q.Claim(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if id != "job-1" {
		t.Fatalf("want job-1, got %s", id)
	}
}

func TestMemoryQueue_ClaimMovesToProcessing(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	ctx := context.Background()

	q.Push(ctx, "job-1")
	q.Claim(ctx)

	processing, err := q.ListProcessing(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(processing) != 1 || processing[0] != "job-1" {
		t.Fatalf("want [job-1] in processing, got %v", processing)
	}
}

func TestMemoryQueue_AckRemovesFromProcessing(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	ctx := context.Background()

	q.Push(ctx, "job-1")
	q.Claim(ctx)
	q.Ack(ctx, "job-1")

	processing, _ := q.ListProcessing(ctx)
	if len(processing) != 0 {
		t.Fatalf("want empty processing list, got %v", processing)
	}
}

func TestMemoryQueue_ClaimBlocksUntilPush(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	ctx := context.Background()

	go func() {
		time.Sleep(20 * time.Millisecond)
		q.Push(ctx, "job-async")
	}()

	id, err := q.Claim(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if id != "job-async" {
		t.Fatalf("want job-async, got %s", id)
	}
}

func TestMemoryQueue_ClaimRespectsContextCancel(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	_, err := q.Claim(ctx)
	if err == nil {
		t.Fatal("want error on cancelled context, got nil")
	}
}

func TestMemoryQueue_MultipleJobsOrdered(t *testing.T) {
	q := queue.NewMemoryQueue(10)
	ctx := context.Background()

	for _, id := range []string{"a", "b", "c"} {
		q.Push(ctx, id)
	}

	for _, want := range []string{"a", "b", "c"} {
		got, err := q.Claim(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("want %s, got %s", want, got)
		}
	}
}
