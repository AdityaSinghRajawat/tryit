package ai

// Response is the provider-agnostic completion output.
type Response struct {
	Content string // raw text — JSON when the request supplied a Schema
	Model   string // model identifier the provider actually used
}
