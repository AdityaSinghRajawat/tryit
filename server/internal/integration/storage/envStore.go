// Package storage holds concrete implementations of the service-owned ports
// (secret.Store, pair.Store, …). Phase 1: envStore (secret) + pairStore.
//
// envStore maps a logical NAME (UPPER_SNAKE) → env var TRYIT_SECRET_<NAME>.
// For basic auth, two env vars: TRYIT_SECRET_<NAME>_USER and _PASS.
// Phase 1 only; replaced by keychainStore in Phase 2.
package storage

import (
	"os"

	"github.com/AdityaSinghRajawat/tryit/server/internal/model"
)

type EnvStore struct{}

func NewEnvStore() *EnvStore { return &EnvStore{} }

// Get implements secret.Store. Returns (Secret, true, nil) if at least one
// matching env var is set; (zero, false, nil) otherwise.
func (e *EnvStore) Get(name string) (model.Secret, bool, error) {
	if name == "" {
		return model.Secret{}, false, nil
	}
	prefix := "TRYIT_SECRET_" + name
	if user := os.Getenv(prefix + "_USER"); user != "" {
		pass := os.Getenv(prefix + "_PASS")
		return model.NewBasicSecret(name, user, pass), true, nil
	}
	if v := os.Getenv(prefix); v != "" {
		return model.NewSecret(name, "bearer", v), true, nil
	}
	return model.Secret{}, false, nil
}
