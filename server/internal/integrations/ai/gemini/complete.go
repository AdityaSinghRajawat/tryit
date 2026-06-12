package gemini

import (
	"context"
	"errors"
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	aiType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/ai"
	"github.com/google/generative-ai-go/genai"
)

// Complete: GenerateContent with system instruction. When a schema is
// supplied we set ResponseMIMEType=application/json and embed the schema in
// the prompt — Gemini's ResponseSchema accepts only the SDK's own Schema
// type, so the caller's JSON-Schema validator enforces shape.
func (p *GeminiProvider) Complete(
	ctx context.Context,
	req aiType.Request,
) (*aiType.Response, error) {
	model := p.client.GenerativeModel(p.model)

	if req.System != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(req.System)},
		}
	}

	temperature := float32(req.Temperature)
	if temperature == 0 {
		temperature = float32(config.GetAITemperature())
	}
	maxTokens := int32(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = int32(config.GetAIMaxTokens())
	}
	model.GenerationConfig.Temperature = &temperature
	model.GenerationConfig.MaxOutputTokens = &maxTokens

	userText := req.User
	if len(req.Schema) > 0 {
		model.GenerationConfig.ResponseMIMEType = "application/json"
		userText = req.User + "\n\nRespond with ONLY a JSON object matching this schema:\n" + string(
			req.Schema,
		)
	}

	resp, err := model.GenerateContent(ctx, genai.Text(userText))
	if err != nil {
		return nil, fmt.Errorf("gemini: %w", err)
	}
	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, errors.New("gemini: empty response")
	}

	var content string
	for _, part := range resp.Candidates[0].Content.Parts {
		if t, ok := part.(genai.Text); ok {
			content += string(t)
		}
	}

	return &aiType.Response{
		Content: content,
		Model:   p.model,
	}, nil
}
