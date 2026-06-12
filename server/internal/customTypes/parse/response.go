package parse

import (
	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
)

type Source string

const (
	SourceCache     Source = "cache"
	SourceSpec      Source = "spec"
	SourceExtractor Source = "extractor"
	SourceProfile   Source = "profile"
	SourceAI        Source = "ai"
)

type Response struct {
	RequestSpec       specType.RequestSpec `json:"requestSpec"`
	Source            Source               `json:"source"`
	Confidence        float64              `json:"confidence"`
	NeedsConfirmation bool                 `json:"needsConfirmation"`
}

func BuildResponse(spec specType.RequestSpec, src Source) *Response {
	return &Response{
		RequestSpec:       spec,
		Source:            src,
		Confidence:        spec.Confidence,
		NeedsConfirmation: spec.Confidence < config.GetParseConfidenceThreshold(),
	}
}
