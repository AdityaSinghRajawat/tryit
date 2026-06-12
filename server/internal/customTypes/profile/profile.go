package profile

import (
	"time"

	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
)

type Source string

const (
	SourceBuiltin Source = "builtin"
	SourceLearned Source = "learned"
)

// SiteProfile pins a host's known request + auth shape; learned overrides builtin.
type SiteProfile struct {
	Host      string            `json:"host"                validate:"required,min=1"`
	Name      string            `json:"name,omitempty"`
	BaseURL   string            `json:"baseUrl"             validate:"required,url"`
	Auth      specType.AuthSpec `json:"auth"`
	Source    Source            `json:"source,omitempty"`
	UpdatedAt time.Time         `json:"updatedAt,omitempty"`
}
