// Package pair owns the pairing logic (§8.1). Token is generated outside the
// service (by the Store at construction time so that it survives restarts);
// Verify performs the constant-time compare and binds the origin on first
// successful match.
package pair

import (
	"crypto/subtle"
	"errors"
	"strings"
)

// Store is the port the pair service depends on. Implementations live under
// internal/integration/storage (Phase 1: pairStore file-backed; Phase 2:
// keychain entry tryit-pair-token).
type Store interface {
	Token() string
	BoundOrigin() string
	SetBoundOrigin(origin string) error
	Reset() error
}

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrInvalidOrigin  = errors.New("invalid origin")
	ErrOriginConflict = errors.New("origin already bound to a different value")
)

type Service struct {
	store Store
}

func New(store Store) *Service { return &Service{store: store} }

// Verify implements §8.1:
//   - If unbound: constant-time compare; on match, bind the calling Origin.
//   - If bound: require Origin == bound AND token match.
//
// Returns the (now-)bound origin on success.
func (s *Service) Verify(token, origin string) (string, error) {
	if subtle.ConstantTimeCompare([]byte(strings.TrimSpace(token)), []byte(s.store.Token())) != 1 {
		return "", ErrInvalidToken
	}
	if !strings.HasPrefix(origin, "chrome-extension://") {
		return "", ErrInvalidOrigin
	}
	bound := s.store.BoundOrigin()
	if bound == "" {
		if err := s.store.SetBoundOrigin(origin); err != nil {
			return "", err
		}
		return origin, nil
	}
	if bound != origin {
		return "", ErrOriginConflict
	}
	return bound, nil
}

func (s *Service) Token() string       { return s.store.Token() }
func (s *Service) BoundOrigin() string { return s.store.BoundOrigin() }
