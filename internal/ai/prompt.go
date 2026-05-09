package ai

import "encoding/json"

// AgentPayload is parsed from job.Payload.
type AgentPayload struct {
	SystemPrompt string            `json:"system_prompt"`
	UserMessage  string            `json:"user_message"`
	Tools        []json.RawMessage `json:"tools"`
	MaxSteps     int               `json:"max_steps"`
}

// AgentResult is written to job.Result on completion.
type AgentResult struct {
	FinalMessage string `json:"final_message"`
	Steps        int    `json:"steps"`
	StopReason   string `json:"stop_reason"`
}
