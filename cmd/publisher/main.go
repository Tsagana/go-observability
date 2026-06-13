package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go-observability/internal/db"
	"go-observability/internal/outbox"
	"go-observability/internal/queue"
	redisclient "go-observability/internal/redis"
)

type config struct {
	databaseURL        string
	redisURL           string
	queuePendingKey    string
	queueProcessingKey string
	queueClaimTimeout  time.Duration
	eventsKey          string
	batchSize          int
	pollInterval       time.Duration
}

func loadConfig() config {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	queuePendingKey := os.Getenv("QUEUE_PENDING_KEY")
	if queuePendingKey == "" {
		queuePendingKey = "jobs:pending"
	}
	queueProcessingKey := os.Getenv("QUEUE_PROCESSING_KEY")
	if queueProcessingKey == "" {
		queueProcessingKey = "jobs:processing"
	}
	queueClaimTimeout, _ := time.ParseDuration(os.Getenv("QUEUE_CLAIM_TIMEOUT"))
	if queueClaimTimeout == 0 {
		queueClaimTimeout = 5 * time.Second
	}
	batchSize, _ := strconv.Atoi(os.Getenv("PUBLISHER_BATCH_SIZE"))
	if batchSize == 0 {
		batchSize = 100
	}
	pollInterval, _ := time.ParseDuration(os.Getenv("PUBLISHER_POLL_INTERVAL"))
	if pollInterval == 0 {
		pollInterval = 1 * time.Second
	}

	eventsKey := os.Getenv("EVENTS_QUEUE_KEY")
	if eventsKey == "" {
		eventsKey = "events:pending"
	}

	return config{
		databaseURL:        os.Getenv("DATABASE_URL"),
		redisURL:           redisURL,
		queuePendingKey:    queuePendingKey,
		queueProcessingKey: queueProcessingKey,
		queueClaimTimeout:  queueClaimTimeout,
		eventsKey:          eventsKey,
		batchSize:          batchSize,
		pollInterval:       pollInterval,
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cfg := loadConfig()

	dbpool, err := db.Connect(context.Background(), cfg.databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	rdb, err := redisclient.Connect(context.Background(), cfg.redisURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()
	slog.Info("redis.connected", "url", cfg.redisURL)

	outboxStore := outbox.NewStore(dbpool)
	jobsQueue := queue.NewRedisQueue(rdb, cfg.queuePendingKey, cfg.queueProcessingKey, cfg.queueClaimTimeout)
	eventsQueue := queue.NewRedisQueue(rdb, cfg.eventsKey, "", 0)
	publisher := outbox.NewPublisher(outboxStore, jobsQueue, eventsQueue, cfg.batchSize, cfg.pollInterval)

	slog.Info("publisher.starting", "batch_size", cfg.batchSize, "poll_interval", cfg.pollInterval)
	if err := publisher.Start(ctx); err != nil {
		log.Fatal(err)
	}
	slog.Info("publisher.shutdown")
}
