package main

import (
	"context"
	"go-observability/internal/api"
	"go-observability/internal/db"
	"go-observability/internal/job"
	"go-observability/internal/worker"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	//Init context
	ctx := context.Background()
	dbpool, err := db.Connect(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	store := job.NewStore(dbpool)
	handler := api.NewHandler(store)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	//Run worker
	w := worker.New(store, 2*time.Second)
	go w.Run(ctx)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
