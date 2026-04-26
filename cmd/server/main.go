package main

import (
	"context"
	"go-observability/internal/api"
	"go-observability/internal/db"
	"go-observability/internal/job"
	"go-observability/internal/worker"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	//Init context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	dbpool, err := db.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	store := job.NewStore(dbpool)
	handler := api.NewHandler(store)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

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

	//Job workers dispatcher
	dispatcher := worker.NewDispatcher(store, workerCount, bufferSize, pollInterval)
	go dispatcher.Run(ctx)

	//Stuck jobs reaper
	interval, _ := time.ParseDuration(os.Getenv("REAPER_INTERVAL"))
	if interval == 0 {
		interval = 60 * time.Second
	}
	stuckJobTimeout, _ := time.ParseDuration(os.Getenv("STUCK_JOB_TIMEOUT"))
	if stuckJobTimeout == 0 {
		stuckJobTimeout = 300 * time.Second
	}

	reaper := worker.NewReaper(store, interval, stuckJobTimeout)
	go reaper.Run(ctx)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
