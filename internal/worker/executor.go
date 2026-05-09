package worker

import (
	"context"
	"encoding/json"
	"go-observability/internal/ai"
	"go-observability/internal/job"
)

func process(ctx context.Context, j *job.Job, client *ai.Client) ([]byte, error) {
	var payload ai.AgentPayload
	if err := json.Unmarshal(j.Payload, &payload); err != nil {
		return nil, job.NewPermanent(err)
	}

	result, err := client.RunAgentLoop(ctx, payload)
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return out, nil

}
