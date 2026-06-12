package gemini

import (
	"context"
	"errors"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	client *genai.Client
	model  string
}

func NewGeminiProvider() (*GeminiProvider, error) {
	apiKey := config.GetGeminiAPIKey()
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY is empty")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &GeminiProvider{
		client: client,
		model:  config.GetDefaultGeminiModel(),
	}, nil
}
