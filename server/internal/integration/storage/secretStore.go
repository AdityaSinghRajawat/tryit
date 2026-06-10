package storage

import (
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/service/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/config"
)

// NewSecretStore returns the configured secret.Store implementation. Phase 1
// only supports "env"; "keychain" and "file" arrive in Phase 2 in this same
// package, preserving the port.
func NewSecretStore(cfg config.Config) (secret.Store, error) {
	switch cfg.SecretsBackend {
	case "env":
		return NewEnvStore(), nil
	default:
		// DECISION (D-P1-5): Phase 1 only knows "env". Surface the misconfig
		// clearly rather than silently falling back.
		return nil, fmt.Errorf("secrets backend %q not implemented in Phase 1 (only \"env\")", cfg.SecretsBackend)
	}
}
