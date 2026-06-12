package parse

type Request struct {
	PageURL             string `json:"pageUrl"                       validate:"required,url"`
	ScopedMarkdown      string `json:"scopedMarkdown,omitempty"`
	AuthSectionMarkdown string `json:"authSectionMarkdown,omitempty"`
	Framework           string `json:"framework,omitempty"`
	// StructuredHint (e.g. OpenAPI operation object) is preferred over markdown.
	StructuredHint any  `json:"structuredHint,omitempty"`
	Force          bool `json:"force,omitempty"`
}
