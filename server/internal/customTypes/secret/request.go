package secret

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// CreateRequest covers all auth types — bearer/apiKey use Value; basic uses
// Username + Password.
type CreateRequest struct {
	Name     string `json:"name"               validate:"required,min=1"`
	Type     string `json:"type"               validate:"required,oneof=bearer basic apiKey"`
	HostHint string `json:"hostHint,omitempty"`
	Value    string `json:"value,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return fmt.Errorf("secret request validation failed: %w", err)
	}
	var nameRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

	if !nameRe.MatchString(r.Name) {
		return errors.New(
			"secret name must be UPPER_SNAKE_CASE (A-Z, 0-9, _; must start with a letter)",
		)
	}
	switch r.Type {
	case "bearer", "apiKey":
		if r.Value == "" {
			return fmt.Errorf("%s secret requires a non-empty value", r.Type)
		}
	case "basic":
		if r.Username == "" || r.Password == "" {
			return errors.New("basic secret requires username and password")
		}
	}
	return nil
}

func (r *CreateRequest) ToRecord() *StoredSecret {
	return &StoredSecret{
		Name:     r.Name,
		Type:     r.Type,
		HostHint: r.HostHint,
		Value:    r.Value,
		User:     r.Username,
		Pass:     r.Password,
	}
}
