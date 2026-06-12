package validations

import "github.com/go-playground/validator/v10"

var structValidator = validator.New(validator.WithRequiredStructEnabled())

func ValidateStruct(s any) error {
	return structValidator.Struct(s)
}
