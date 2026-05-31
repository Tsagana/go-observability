package main

import (
	"context"
	"go-observability/internal/ai"
	"go-observability/internal/db"
	"go-observability/internal/job"
	"go-observability/internal/queue"
	redisclient "go-observability/internal/redis"
	"go-observability/internal/worker"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type config struct {
	databaseURL        string
	redisURL           string
	queuePendingKey    string
	queueProcessingKey string
	queueClaimTimeout  time.Duration

	workerCount  int
	bufferSize   int
	pollInterval time.Duration

	reaperInterval  time.Duration
	stuckJobTimeout time.Duration

	anthropicAPIKey   string
	aiStepTimeout     time.Duration
	aiJobTimeout      time.Duration
	aiMaxRetries      int
	aiDefaultMaxSteps int
}

func loadConfig() config {
	workerCount, _ := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	if workerCount == 0 {
		workerCount = 5
	}
	bufferSize, _ := strconv.Atoi(os.Getenv("JOB_CHANNEL_BUFFER"))
	if bufferSize == 0 {
		bufferSize = workerCount
	}
	pollInterval, _ := time.ParseDuration(os.Getenv("POLL_INTERVAL"))
	if pollInterval == 0 {
		pollInterval = 2 * time.Second
	}
	reaperInterval, _ := time.ParseDuration(os.Getenv("REAPER_INTERVAL"))
	if reaperInterval == 0 {
		reaperInterval = 60 * time.Second
	}
	stuckJobTimeout, _ := time.ParseDuration(os.Getenv("STUCK_JOB_TIMEOUT"))
	if stuckJobTimeout == 0 {
		stuckJobTimeout = 300 * time.Second
	}
	aiStepTimeout, _ := time.ParseDuration(os.Getenv("AI_STEP_TIMEOUT"))
	if aiStepTimeout == 0 {
		aiStepTimeout = 30 * time.Second
	}
	aiJobTimeout, _ := time.ParseDuration(os.Getenv("AI_JOB_TIMEOUT"))
	if aiJobTimeout == 0 {
		aiJobTimeout = 5 * time.Minute
	}
	aiMaxRetries, _ := strconv.Atoi(os.Getenv("AI_MAX_RETRIES"))
	if aiMaxRetries == 0 {
		aiMaxRetries = 3
	}
	aiDefaultMaxSteps, _ := strconv.Atoi(os.Getenv("AI_DEFAULT_MAX_STEPS"))
	if aiDefaultMaxSteps == 0 {
		aiDefaultMaxSteps = 25
	}
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

	return config{
		databaseURL:        os.Getenv("DATABASE_URL"),
		redisURL:           redisURL,
		queuePendingKey:    queuePendingKey,
		queueProcessingKey: queueProcessingKey,
		queueClaimTimeout:  queueClaimTimeout,
		workerCount:        workerCount,
		bufferSize:         bufferSize,
		pollInterval:       pollInterval,
		reaperInterval:     reaperInterval,
		stuckJobTimeout:    stuckJobTimeout,
		anthropicAPIKey:    os.Getenv("ANTHROPIC_API_KEY"),
		aiStepTimeout:      aiStepTimeout,
		aiJobTimeout:       aiJobTimeout,
		aiMaxRetries:       aiMaxRetries,
		aiDefaultMaxSteps:  aiDefaultMaxSteps,
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cfg := loadConfig()
	if cfg.anthropicAPIKey == "" {
		log.Fatal("ANTHROPIC_API_KEY is required")
	}

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

	store := job.NewStore(dbpool)
	q := queue.NewRedisQueue(rdb, cfg.queuePendingKey, cfg.queueProcessingKey, cfg.queueClaimTimeout)
	aiClient := ai.NewClient(cfg.anthropicAPIKey, cfg.aiStepTimeout, cfg.aiJobTimeout, cfg.aiMaxRetries, cfg.aiDefaultMaxSteps)
	dispatcher := worker.NewDispatcher(store, cfg.workerCount, cfg.bufferSize, q, aiClient)
	reaper := worker.NewReaper(store, q, cfg.reaperInterval, cfg.stuckJobTimeout)

	go reaper.Run(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatcher.Run(ctx)
	}()

	slog.Info("worker.starting", "workers", cfg.workerCount)
	wg.Wait()
	slog.Info("worker.shutdown")
}
