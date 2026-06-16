// Package secret resolves and persists named credentials. The active store
// is chosen at boot via TRYIT_SECRETS_BACKEND.
package secret

import (
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/services/secret/stores"
)

// SecretStore is the persistence seam. Get returns (nil, nil) for "not
// found"; Delete reports existence through the deleted bool.
type SecretStore interface {
	Get(name string) (*secretType.StoredSecret, error)
	List() ([]*secretType.StoredSecret, error)
	Create(s *secretType.StoredSecret) error
	Delete(name string) (deleted bool, err error)
}

type SecretService struct {
	store SecretStore
}

func NewSecretService(store SecretStore) *SecretService {
	return &SecretService{store: store}
}

func NewSecretStore() (SecretStore, error) {
	switch config.GetSecretsBackend() {
	case "", config.GetSecretsProviderEnv():
		return stores.NewEnvStore(), nil
	case config.GetSecretsProviderKeychain():
		return stores.NewKeychainStore(), nil
	case config.GetSecretsProviderFile():
		return stores.NewFileStore(config.GetSecretsFile(), config.GetSecretsPassphrase())
	default:
		return nil, fmt.Errorf(
			"unknown TRYIT_SECRETS_BACKEND %q (want: env | keychain | file)",
			config.GetSecretsBackend(),
		)
	}
}
