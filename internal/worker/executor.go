package worker

import (
	"log"
	"time"

	"go-observability/internal/job"
)

func process(j *job.Job) ([]byte, error) {
	time.Sleep(15 * time.Second) // simulate slow job
	log.Printf("processing job %s", j.ID)
	return []byte(`{"result": true}`), nil
}
