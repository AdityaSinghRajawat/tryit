package parse

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Request struct {
	PageURL             string `json:"pageUrl"                       validate:"required,url"`
	ScopedMarkdown      string `json:"scopedMarkdown,omitempty"`
	AuthSectionMarkdown string `json:"authSectionMarkdown,omitempty"`
	Framework           string `json:"framework,omitempty"`
	// StructuredHint (e.g. OpenAPI operation object) is preferred over markdown.
	StructuredHint any  `json:"structuredHint,omitempty"`
	Force          bool `json:"force,omitempty"`
}

func (r *Request) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return fmt.Errorf("parse request validation failed: %w", err)
	}
	return nil
}
