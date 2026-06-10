// pairStore persists the pairing token + bound origin to a 0600 file at
// ~/.tryit/pair.json (DECISION D-P1-3 — Phase 1 file-backed; Phase 2 swaps to
// keychain entry "tryit-pair-token" behind the same pair.Store port).
//
// On first construction the file does not exist: the store generates a fresh
// token and writes it before returning. The boundOrigin field is "" until the
// first successful /pair call calls SetBoundOrigin.
package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/AdityaSinghRajawat/tryit/server/internal/helpers"
)

type pairFile struct {
	Token       string `json:"token"`
	BoundOrigin string `json:"boundOrigin"`
}

type PairStore struct {
	path string
	mu   sync.RWMutex
	data pairFile
}

func NewPairStore(path string) (*PairStore, bool, error) {
	if path == "" {
		return nil, false, errors.New("pair store path is empty (set TRYIT_PAIR_FILE or $HOME)")
	}
	s := &PairStore{path: path}
	freshlyGenerated, err := s.load()
	if err != nil {
		return nil, false, err
	}
	return s, freshlyGenerated, nil
}

func (s *PairStore) load() (bool, error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return false, err
	}
	b, err := os.ReadFile(s.path)
	if err == nil {
		var f pairFile
		if jerr := json.Unmarshal(b, &f); jerr != nil {
			return false, jerr
		}
		s.data = f
		if f.Token == "" {
			// Corrupt-ish: re-generate.
			return s.regenerate()
		}
		return false, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	return s.regenerate()
}

func (s *PairStore) regenerate() (bool, error) {
	tok, err := helpers.Token()
	if err != nil {
		return false, err
	}
	s.data = pairFile{Token: tok}
	return true, s.flush()
}

func (s *PairStore) flush() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o600)
}

func (s *PairStore) Token() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.Token
}

func (s *PairStore) BoundOrigin() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.BoundOrigin
}

func (s *PairStore) SetBoundOrigin(origin string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.BoundOrigin = origin
	return s.flush()
}

func (s *PairStore) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.regenerate()
	return err
}
