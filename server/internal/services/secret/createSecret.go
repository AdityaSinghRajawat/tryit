package secret

import (
	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
)

func (s *SecretService) CreateSecret(record *secretType.StoredSecret) *config.CustomError {
	if err := s.store.Create(record); err != nil {
		return config.NewCustomError(err, config.GetErrCodeInternal())
	}

	return nil
}
