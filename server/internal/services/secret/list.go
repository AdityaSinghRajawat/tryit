package secret

import (
	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
)

func (s *SecretService) ListSecrets() ([]secretType.Info, *config.CustomError) {
	stored, err := s.store.List()
	if err != nil {
		return nil, config.NewCustomError(err, config.GetErrCodeInternal())
	}

	out := make([]secretType.Info, 0, len(stored))
	for _, rec := range stored {
		out = append(out, secretType.Info{
			Name:     rec.Name,
			Type:     rec.Type,
			HostHint: rec.HostHint,
		})
	}

	return out, nil
}
