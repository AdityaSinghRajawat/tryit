package execute

import (
	"fmt"

	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
)

type Request struct {
	RequestSpec specType.RequestSpec `json:"requestSpec"          validate:"required"`
	SecretRefs  map[string]string    `json:"secretRefs,omitempty"`
}

func (r *Request) Validate() error {
	if err := r.RequestSpec.Validate(); err != nil {
		return fmt.Errorf("execute request validation failed: %w", err)
	}
	return nil
}
