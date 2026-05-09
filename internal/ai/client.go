package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-observability/internal/job"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Client struct {
	anthropic       *anthropic.Client
	stepTimeout     time.Duration
	jobTimeout      time.Duration
	maxRetries      int
	defaultMaxSteps int
}

func NewClient(apiKey string, stepTimeout, jobTimeout time.Duration, maxRetries, defaultMaxSteps int) *Client {
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Client{
		anthropic:       &c,
		stepTimeout:     stepTimeout,
		jobTimeout:      jobTimeout,
		maxRetries:      maxRetries,
		defaultMaxSteps: defaultMaxSteps,
	}
}

func (c *Client) RunAgentLoop(ctx context.Context, payload AgentPayload) (AgentResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.jobTimeout)
	defer cancel()

	maxSteps := payload.MaxSteps
	if maxSteps == 0 {
		maxSteps = c.defaultMaxSteps
	}

	tools := make([]anthropic.ToolUnionParam, len(payload.Tools))
	for i, raw := range payload.Tools {
		var t anthropic.ToolParam
		if err := json.Unmarshal(raw, &t); err != nil {
			return AgentResult{}, job.NewPermanent(fmt.Errorf("invalid tool schema at index %d: %w", i, err))
		}
		tc := t // copy before taking address
		tools[i] = anthropic.ToolUnionParam{OfTool: &tc}
	}

	//build initial messages slice from payload.SystemPrompt + payload.UserMessage
	params := anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_6,
		MaxTokens: 1024,
		System:    []anthropic.TextBlockParam{{Text: payload.SystemPrompt}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(payload.UserMessage)),
		},
		Tools: tools,
	}

	stepCount := 0
	for stepCount < maxSteps {
		stepCtx, stepCancel := context.WithTimeout(ctx, c.stepTimeout)

		// retry loop for transient API errors
		var resp *anthropic.Message
		var err error
		for attempt := 0; attempt <= c.maxRetries; attempt++ {
			resp, err = c.anthropic.Messages.New(stepCtx, params)
			if err == nil {
				break
			}
			if !job.IsRetryable(err) {
				stepCancel()
				return AgentResult{}, job.NewPermanent(err)
			}
			if attempt == c.maxRetries {
				stepCancel()
				return AgentResult{}, err // retries exhausted, dispatcher retries the job
			}
			time.Sleep(time.Duration(1<<attempt) * time.Second) // exponential backoff
		}

		params.Messages = append(params.Messages, resp.ToParam())
		stepCount++

		switch resp.StopReason {
		case anthropic.StopReasonEndTurn:
			stepCancel()
			//extract text from resp.Content using block.AsText()
			var finalMessage string
			for _, block := range resp.Content {
				if block.Type == "text" {
					finalMessage += block.AsText().Text
				}
			}

			return AgentResult{FinalMessage: finalMessage, Steps: stepCount, StopReason: string(resp.StopReason)}, nil

		case anthropic.StopReasonToolUse:
			var toolResults []anthropic.ContentBlockParamUnion
			for _, block := range resp.Content {
				if block.Type != "tool_use" {
					continue
				}
				toolUse := block.AsToolUse()
				result, isError := executeTool(toolUse.Name, toolUse.Input)
				toolResults = append(toolResults, anthropic.NewToolResultBlock(toolUse.ID, result, isError))
			}
			params.Messages = append(params.Messages, anthropic.NewUserMessage(toolResults...))

		}

		stepCancel()
	}

	return AgentResult{}, job.NewPermanent(fmt.Errorf("max steps exceeded: %d", maxSteps))
}

func executeTool(name string, input json.RawMessage) (string, bool) {
	switch name {
	case "calculator":
		return runCalculator(input)
	default:
		return fmt.Sprintf("unknown tool: %s", name), true
	}
}

func runCalculator(input json.RawMessage) (string, bool) {
	var params struct {
		Expression string `json:"expression"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "invalid input: " + err.Error(), true
	}

	var a, b float64
	var op string
	if _, err := fmt.Sscanf(params.Expression, "%f %s %f", &a, &op, &b); err != nil {
		return "invalid expression: " + err.Error(), true
	}

	switch op {
	case "+":
		return fmt.Sprintf("%g", a+b), false
	case "-":
		return fmt.Sprintf("%g", a-b), false
	case "*":
		return fmt.Sprintf("%g", a*b), false
	case "/":
		if b == 0 {
			return "division by zero", true
		}
		return fmt.Sprintf("%g", a/b), false
	default:
		return "unsupported operator: " + op, true
	}
}
