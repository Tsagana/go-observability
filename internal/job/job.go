package job

import (
	"encoding/json"
	"math"
	"time"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

const MaxRetry = 3

type Job struct {
	ID         string          `json:"id"`
	Status     Status          `json:"status"`
	Payload    json.RawMessage `json:"payload"`
	Result     json.RawMessage `json:"result,omitempty"`
	Error      *string         `json:"error,omitempty"`
	RetryCount int             `json:"retry_count"`
	RetryAfter *time.Time      `json:"retry_after,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func calculateBackoff(retryCnt int) int {
	return int(math.Pow(2, float64(retryCnt)))
}
