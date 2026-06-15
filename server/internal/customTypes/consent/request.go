package consent

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Request struct {
	Secret string `json:"secret" validate:"required,min=1"`
	Host   string `json:"host"   validate:"required,min=1"`
}

func (r *Request) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return fmt.Errorf("consent request validation failed: %w", err)
	}
	return nil
}
