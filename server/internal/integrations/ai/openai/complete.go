package openai

import (
	"context"
	"errors"
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	aiType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/ai"
	gopenai "github.com/sashabaranov/go-openai"
)

// Complete: chat-completion with optional JSON-object mode when a schema is
// supplied. We use JSON-object rather than the strict JSON-schema mode so the
// caller's schema (which may use features OpenAI's strict mode rejects) is
// validated by our own validator instead.
func (p *OpenAIProvider) Complete(
	ctx context.Context,
	req aiType.Request,
) (*aiType.Response, error) {
	msgs := []gopenai.ChatCompletionMessage{
		{Role: gopenai.ChatMessageRoleSystem, Content: req.System},
		{Role: gopenai.ChatMessageRoleUser, Content: req.User},
	}

	chatReq := gopenai.ChatCompletionRequest{
		Model:       p.model,
		Messages:    msgs,
		Temperature: float32(req.Temperature),
		MaxTokens:   req.MaxTokens,
	}
	if req.Temperature == 0 && chatReq.Temperature == 0 {
		chatReq.Temperature = float32(config.GetAITemperature())
	}
	if req.MaxTokens == 0 {
		chatReq.MaxTokens = config.GetAIMaxTokens()
	}
	if len(req.Schema) > 0 {
		chatReq.ResponseFormat = &gopenai.ChatCompletionResponseFormat{
			Type: gopenai.ChatCompletionResponseFormatTypeJSONObject,
		}
	}

	resp, err := p.client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("openai: empty response")
	}

	return &aiType.Response{
		Content: resp.Choices[0].Message.Content,
		Model:   p.model,
	}, nil
}
