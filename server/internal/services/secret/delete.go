package secret

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

func (s *SecretService) DeleteSecret(name string) *config.CustomError {
	if strings.TrimSpace(name) == "" {
		return config.NewCustomError(
			errors.New("secret name is empty"),
			config.GetErrCodeInvalidRequest(),
		)
	}

	deleted, err := s.store.Delete(name)
	if err != nil {
		return config.NewCustomError(err, config.GetErrCodeInternal())
	}

	if !deleted {
		return config.NewCustomError(
			fmt.Errorf("secret %q not found", name),
			config.GetErrCodeSecretNotFound(),
		)
	}

	return nil
}
