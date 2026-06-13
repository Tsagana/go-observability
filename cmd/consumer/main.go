package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-observability/internal/db"
	"go-observability/internal/events"
	"go-observability/internal/queue"
	redisclient "go-observability/internal/redis"
)

type config struct {
	databaseURL  string
	redisURL     string
	eventsKey    string
	blockTimeout time.Duration
}

func loadConfig() config {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	eventsKey := os.Getenv("EVENTS_QUEUE_KEY")
	if eventsKey == "" {
		eventsKey = "events:pending"
	}

	blockTimeout, _ := time.ParseDuration(os.Getenv("CONSUMER_BLOCK_TIMEOUT"))
	if blockTimeout == 0 {
		blockTimeout = 5 * time.Second
	}

	return config{
		databaseURL:  os.Getenv("DATABASE_URL"),
		redisURL:     redisURL,
		eventsKey:    eventsKey,
		blockTimeout: blockTimeout,
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

	eventsStore := events.NewStore(dbpool)
	eventsQueue := queue.NewRedisQueue(rdb, cfg.eventsKey, "", cfg.blockTimeout)
	consumer := events.NewConsumer(eventsStore, eventsQueue)

	slog.Info("consumer.starting", "events_key", cfg.eventsKey, "block_timeout", cfg.blockTimeout)
	if err := consumer.Start(ctx); err != nil {
		log.Fatal(err)
	}
	slog.Info("consumer.shutdown")
}
