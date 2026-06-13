package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-observability/internal/api"
	"go-observability/internal/db"
	"go-observability/internal/job"
	"go-observability/internal/outbox"
)

type config struct {
	databaseURL string
	port        string
}

func loadConfig() config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return config{
		databaseURL: os.Getenv("DATABASE_URL"),
		port:        port,
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

	store := job.NewStore(dbpool)
	outboxStore := outbox.NewStore(dbpool)

	handler := api.NewHandler(store, outboxStore, dbpool)
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
