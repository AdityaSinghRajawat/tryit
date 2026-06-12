package ai

// Request is the provider-agnostic completion input.
type Request struct {
	System      string  // system prompt
	User        string  // user message
	Schema      []byte  // optional JSON Schema for structured output (provider may embed in prompt)
	Temperature float64 // 0 for deterministic
	MaxTokens   int     // upper bound on response tokens
}
