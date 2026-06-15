// Package consent persists per-(secret, host) grants (IMPL §8.4). Execute
// checks Granted before injecting a secret; the panel calls Grant after the
// ConsentDialog approval.
package consent

import (
	"errors"
	"sync"

	consentType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/consent"
)

type ConsentService struct {
	path    string
	mu      sync.RWMutex
	records []consentType.Record
}

func NewConsentService(path string) (*ConsentService, error) {
	if path == "" {
		return nil, errors.New("consent file path is empty (set TRYIT_CONSENT_FILE or $HOME)")
	}

	s := &ConsentService{path: path}
	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}
