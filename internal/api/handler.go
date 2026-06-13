package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-observability/internal/job"
	"go-observability/internal/outbox"
)

type Handler struct {
	store       *job.Store
	outboxStore *outbox.Store
	pool        *pgxpool.Pool
}

func NewHandler(store *job.Store, outboxStore *outbox.Store, pool *pgxpool.Pool) *Handler {
	return &Handler{store: store, outboxStore: outboxStore, pool: pool}
}

type createJobRequest struct {
	Payload json.RawMessage `json:"payload"`
}

type createJobResponse struct {
	ID     string     `json:"id"`
	Status job.Status `json:"status"`
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		slog.Error("pool.begin failed", "error", err)
		http.Error(w, "error when creating job", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	createdJob, err := h.store.CreateTx(ctx, tx, req.Payload)
	if err != nil {
		slog.Error("store.createTx failed", "error", err)
		http.Error(w, "error when creating job", http.StatusInternalServerError)
		return
	}

	eventPayload, _ := json.Marshal(map[string]string{"job_id": createdJob.ID})

	if err := h.outboxStore.InsertTx(ctx, tx, "job.created", eventPayload); err != nil {
		slog.Error("outbox.insertTx failed", "error", err)
		http.Error(w, "error when creating job", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("tx.commit failed", "error", err)
		http.Error(w, "error when creating job", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusAccepted, createJobResponse{ID: createdJob.ID, Status: createdJob.Status})
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.store.Get(r.Context(), id)
	if err != nil {
		slog.Error("store.get failed", "error", err)
		http.Error(w, "error when getting job", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, createJobResponse{ID: job.ID, Status: job.Status})

}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok health here"))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
