package generate

import (
	"fmt"

	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
	"github.com/go-playground/validator/v10"
)

type Language string

const (
	LanguageCurl       Language = "curl"
	LanguagePython     Language = "python"
	LanguageJavaScript Language = "javascript"
	LanguageTypeScript Language = "typescript"
	LanguageGo         Language = "go"
)

type GenerateRequest struct {
	RequestSpec specType.RequestSpec `json:"requestSpec"         validate:"required"`
	Language    Language             `json:"language"            validate:"required,oneof=curl python javascript typescript go"`
	Idiomatic   bool                 `json:"idiomatic,omitempty"`
}

func (r *GenerateRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return fmt.Errorf("generate request validation failed: %w", err)
	}
	if err := r.RequestSpec.Validate(); err != nil {
		return fmt.Errorf("requestSpec invalid: %w", err)
	}
	return nil
}
