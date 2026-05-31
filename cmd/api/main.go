package main

import (
	"context"
	"fmt"
	"go-observability/internal/api"
	"go-observability/internal/db"
	"go-observability/internal/job"
	"go-observability/internal/queue"
	redisclient "go-observability/internal/redis"
	"log"
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"time"
	"os/signal"
)

type config struct {
	databaseURL        string
	redisURL           string
	queuePendingKey    string
	queueProcessingKey string
	queueClaimTimeout  time.Duration
	port               string
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return config{
		databaseURL:        os.Getenv("DATABASE_URL"),
		redisURL:           redisURL,
		queuePendingKey:    queuePendingKey,
		queueProcessingKey: queueProcessingKey,
		queueClaimTimeout:  queueClaimTimeout,
		port:               port,
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

	store := job.NewStore(dbpool)
	q := queue.NewRedisQueue(rdb, cfg.queuePendingKey, cfg.queueProcessingKey, cfg.queueClaimTimeout)

	handler := api.NewHandler(store, q)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	addr := fmt.Sprintf(":%s", cfg.port)
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	slog.Info("api.starting", "addr", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	slog.Info("api.shutdown")
}
