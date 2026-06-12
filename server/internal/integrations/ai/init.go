// Package ai is the Strategy + Factory layer for AI providers. The active
// provider is chosen once at startup from AI_PROVIDER; adding a new provider
// is one new sub-package + one case in NewAIProvider.
package ai

import (
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai/anthropic"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai/gemini"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai/ollama"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai/openai"
)

func NewAIProvider() (AIProvider, error) {
	switch config.GetAIProvider() {
	case config.GetAIProviderAnthropic():
		return anthropic.NewAnthropicProvider()
	case config.GetAIProviderOpenAI():
		return openai.NewOpenAIProvider()
	case config.GetAIProviderGemini():
		return gemini.NewGeminiProvider()
	case config.GetAIProviderOllama():
		return ollama.NewOllamaProvider()
	default:
		return nil, fmt.Errorf(
			"unknown AI_PROVIDER %q (want: anthropic | openai | gemini | ollama)",
			config.GetAIProvider(),
		)
	}
}
