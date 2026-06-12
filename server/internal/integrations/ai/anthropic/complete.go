package anthropic

import (
	"context"
	"errors"
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	aiType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/ai"
	anth "github.com/anthropics/anthropic-sdk-go"
)

// Complete: Messages API with system + user message. When a schema is
// supplied, we ask the model to emit JSON only — the caller's validator
// enforces shape.
func (p *AnthropicProvider) Complete(
	ctx context.Context,
	req aiType.Request,
) (*aiType.Response, error) {
	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = int64(config.GetAIMaxTokens())
	}
	temperature := req.Temperature
	if temperature == 0 {
		temperature = config.GetAITemperature()
	}

	userText := req.User
	if len(req.Schema) > 0 {
		userText = req.User + "\n\nRespond with ONLY a JSON object matching this schema (no prose, no fences):\n" + string(
			req.Schema,
		)
	}

	params := anth.MessageNewParams{
		Model:       anth.Model(p.model),
		MaxTokens:   maxTokens,
		Temperature: anth.Float(temperature),
		System: []anth.TextBlockParam{
			{Text: req.System},
		},
		Messages: []anth.MessageParam{
			anth.NewUserMessage(anth.NewTextBlock(userText)),
		},
	}

	msg, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("anthropic: %w", err)
	}
	if len(msg.Content) == 0 {
		return nil, errors.New("anthropic: empty response")
	}

	var content string
	for _, block := range msg.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &aiType.Response{
		Content: content,
		Model:   p.model,
	}, nil
}
