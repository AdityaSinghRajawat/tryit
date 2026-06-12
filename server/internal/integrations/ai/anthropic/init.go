package anthropic

import (
	"errors"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	anth "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	client *anth.Client
	model  string
}

func NewAnthropicProvider() (*AnthropicProvider, error) {
	apiKey := config.GetAnthropicAPIKey()
	if apiKey == "" {
		return nil, errors.New("ANTHROPIC_API_KEY is empty")
	}

	client := anth.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicProvider{
		client: &client,
		model:  config.GetDefaultAnthropicModel(),
	}, nil
}
