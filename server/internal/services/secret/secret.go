package secret

import (
	"errors"
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
)

func (s *SecretService) Resolve(name string) (secretType.Secret, *config.CustomError) {
	if name == "" {
		return secretType.Secret{}, config.NewCustomError(
			errors.New("secret name is empty"),
			config.GetErrCodeInvalidRequest(),
		)
	}

	prefix := config.GetSecretEnvPrefix() + name
	if user := config.GetEnvByKey(prefix + config.GetSecretUserSuffix()); user != "" {
		pass := config.GetEnvByKey(prefix + config.GetSecretPassSuffix())
		return secretType.NewBasic(name, user, pass), nil
	}

	if v := config.GetEnvByKey(prefix); v != "" {
		return secretType.New(name, "bearer", v), nil
	}

	return secretType.Secret{}, config.NewCustomError(
		fmt.Errorf("secret %q not found", name),
		config.GetErrCodeSecretNotFound(),
	)
}
