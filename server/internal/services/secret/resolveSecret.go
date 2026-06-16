package secret

import (
	"errors"
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
)

func (s *SecretService) ResolveSecret(name string) (secretType.Secret, *config.CustomError) {
	if name == "" {
		return secretType.Secret{}, config.NewCustomError(
			errors.New("secret name is empty"),
			config.GetErrCodeInvalidRequest(),
		)
	}

	rec, err := s.store.Get(name)
	if err != nil {
		return secretType.Secret{}, config.NewCustomError(err, config.GetErrCodeInternal())
	}

	if rec == nil {
		return secretType.Secret{}, config.NewCustomError(
			fmt.Errorf("secret %q not found", name),
			config.GetErrCodeSecretNotFound(),
		)
	}

	return recordToSecret(rec), nil
}

func recordToSecret(rec *secretType.StoredSecret) secretType.Secret {
	if rec.Type == "basic" {
		return secretType.NewBasic(rec.Name, rec.User, rec.Pass)
	}

	return secretType.New(rec.Name, rec.Type, rec.Value)
}
