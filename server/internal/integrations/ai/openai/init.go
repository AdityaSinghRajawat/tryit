package openai

import (
	"errors"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	gopenai "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *gopenai.Client
	model  string
}

func NewOpenAIProvider() (*OpenAIProvider, error) {
	apiKey := config.GetOpenAIAPIKey()
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY is empty")
	}

	return &OpenAIProvider{
		client: gopenai.NewClient(apiKey),
		model:  config.GetDefaultOpenAIModel(),
	}, nil
}
