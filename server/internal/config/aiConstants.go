package config

type aiConsts struct {
	providerAnthropic string
	providerOpenAI    string
	providerGemini    string
	providerOllama    string

	defaultAnthropicModel string
	defaultOpenAIModel    string
	defaultGeminiModel    string
	defaultOllamaModel    string

	defaultOllamaHost string

	temperature float64
	maxTokens   int
}

var aiI = &aiConsts{
	providerAnthropic: "anthropic",
	providerOpenAI:    "openai",
	providerGemini:    "gemini",
	providerOllama:    "ollama",

	defaultAnthropicModel: "claude-sonnet-4-5-20250929",
	defaultOpenAIModel:    "gpt-4o-mini",
	defaultGeminiModel:    "gemini-1.5-flash",
	defaultOllamaModel:    "llama3",

	defaultOllamaHost: "http://localhost:11434",

	temperature: 0,
	maxTokens:   1500,
}

func GetAIProviderAnthropic() string { return aiI.providerAnthropic }
func GetAIProviderOpenAI() string    { return aiI.providerOpenAI }
func GetAIProviderGemini() string    { return aiI.providerGemini }
func GetAIProviderOllama() string    { return aiI.providerOllama }

func GetDefaultAnthropicModel() string { return aiI.defaultAnthropicModel }
func GetDefaultOpenAIModel() string    { return aiI.defaultOpenAIModel }
func GetDefaultGeminiModel() string    { return aiI.defaultGeminiModel }
func GetDefaultOllamaModel() string    { return aiI.defaultOllamaModel }

func GetDefaultOllamaHost() string { return aiI.defaultOllamaHost }

func GetAITemperature() float64 { return aiI.temperature }
func GetAIMaxTokens() int       { return aiI.maxTokens }
