package ollama

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	aiType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/ai"
	"github.com/ollama/ollama/api"
)

// Complete: streaming-disabled Ollama chat call. Connection failures (host
// unreachable, model not pulled, etc.) are wrapped with the configured host
// so the operator sees what was being attempted.
func (p *OllamaProvider) Complete(
	ctx context.Context,
	req aiType.Request,
) (*aiType.Response, error) {
	temperature := req.Temperature
	if temperature == 0 {
		temperature = config.GetAITemperature()
	}
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = config.GetAIMaxTokens()
	}

	messages := []api.Message{
		{Role: "system", Content: req.System},
		{Role: "user", Content: req.User},
	}

	stream := false
	chatReq := &api.ChatRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   &stream,
		Options: map[string]any{
			"temperature": temperature,
			"num_predict": maxTokens,
		},
	}
	if len(req.Schema) > 0 {
		chatReq.Format = json.RawMessage(`"json"`)
	}

	var content string
	err := p.client.Chat(ctx, chatReq, func(resp api.ChatResponse) error {
		content += resp.Message.Content
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ollama (%s): %w", p.host, err)
	}

	return &aiType.Response{
		Content: content,
		Model:   p.model,
	}, nil
}
