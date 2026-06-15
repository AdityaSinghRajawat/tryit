package pair

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Request struct {
	Token string `json:"token" validate:"required,min=1"`
}

func (r *Request) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return fmt.Errorf("pair request validation failed: %w", err)
	}
	return nil
}
