// Package generate renders a RequestSpec into a runnable code snippet
// (curl / python / javascript). Secrets surface as env-var reads — the
// snippet is meant to be copy-pasted into a shell or editor and run.
package generate

import (
	"text/template"

	"github.com/AdityaSinghRajawat/tryit/server/internal/templates/code"
)

type CodegenService struct {
	templates *template.Template
}

func NewCodegenService() (*CodegenService, error) {
	t, err := template.ParseFS(code.TemplatesFS, "*.tmpl")
	if err != nil {
		return nil, err
	}

	return &CodegenService{templates: t}, nil
}
