package consent

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	consentType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/consent"
)

func (s *ConsentService) load() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(b) == 0 {
		return nil
	}

	return json.Unmarshal(b, &s.records)
}

func (s *ConsentService) flush() error {
	b, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, b, 0o600)
}

func (s *ConsentService) findIndex(secret, host string) int {
	for i, r := range s.records {
		if r.Secret == secret && r.Host == host {
			return i
		}
	}

	return -1
}

func (s *ConsentService) ListConsents() []consentType.Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]consentType.Record, len(s.records))
	copy(out, s.records)

	return out
}
