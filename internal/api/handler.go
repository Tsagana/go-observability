package api

import (
	"encoding/json"
	"log"
	"net/http"

	"go-observability/internal/job"
)

type Handler struct {
	store *job.Store
}

func NewHandler(store *job.Store) *Handler {
	return &Handler{store: store}
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

	createdJob, err := h.store.Create(r.Context(), req.Payload)
	if err != nil {
		log.Println("store.Create error:", err)
		http.Error(w, "error when creating job", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusAccepted, createJobResponse{ID: createdJob.ID, Status: createdJob.Status})
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.store.Get(r.Context(), id)
	if err != nil {
		log.Println("store.Get error:", err)
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
