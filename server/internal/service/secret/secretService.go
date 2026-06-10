// Package secret owns the Store port (services own their ports — IMPL §1).
// Phase 1 surface is just Resolve, used by service/execute. Phase 2 will add
// list/create/delete + consent management + profile learning, still behind
// the same package.
package secret

import (
	"errors"

	"github.com/AdityaSinghRajawat/tryit/server/internal/model"
)

// Store is the port the secret service depends on. Implementations live under
// internal/integration/storage (envStore in Phase 1; keychainStore + fileStore
// in Phase 2).
type Store interface {
	// Get returns the stored secret by logical NAME (UPPER_SNAKE). The bool
	// is false (with nil err) when no such secret exists.
	Get(name string) (model.Secret, bool, error)
}

var ErrNotFound = errors.New("secret not found")

type Service struct {
	store Store
}

func New(store Store) *Service { return &Service{store: store} }

// Resolve returns the stored secret by NAME or ErrNotFound.
func (s *Service) Resolve(name string) (model.Secret, error) {
	v, ok, err := s.store.Get(name)
	if err != nil {
		return model.Secret{}, err
	}
	if !ok {
		return model.Secret{}, ErrNotFound
	}
	return v, nil
}
