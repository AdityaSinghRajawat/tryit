package stores

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// FileStore backs SecretStore with a NaCl-secretbox-encrypted JSON file
// keyed by a scrypt-derived symmetric key. File layout: nonce || ciphertext.
type FileStore struct {
	path       string
	passphrase string

	mu     sync.Mutex
	loaded bool
	cache  []*secretType.StoredSecret
}

type fileVault struct {
	Secrets []*secretType.StoredSecret `json:"secrets"`
}

func NewFileStore(path, passphrase string) (*FileStore, error) {
	if path == "" || passphrase == "" {
		return nil, errors.New(
			"file backend requires TRYIT_SECRETS_FILE and TRYIT_SECRETS_PASSPHRASE",
		)
	}

	return &FileStore{path: path, passphrase: passphrase}, nil
}

func (s *FileStore) Get(name string) (*secretType.StoredSecret, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, err
	}

	if i := s.findIndex(name); i >= 0 {
		return cloneRecord(s.cache[i]), nil
	}

	return nil, nil
}

func (s *FileStore) List() ([]*secretType.StoredSecret, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, err
	}

	out := make([]*secretType.StoredSecret, len(s.cache))
	for i, rec := range s.cache {
		out[i] = cloneRecord(rec)
	}

	return out, nil
}

func (s *FileStore) Create(rec *secretType.StoredSecret) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {

		return err
	}
	if i := s.findIndex(rec.Name); i >= 0 {
		s.cache[i] = rec
	} else {
		s.cache = append(s.cache, rec)
	}

	return s.flushLocked()
}

func (s *FileStore) Delete(name string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return false, err
	}

	i := s.findIndex(name)
	if i < 0 {
		return false, nil
	}

	s.cache = append(s.cache[:i], s.cache[i+1:]...)
	return true, s.flushLocked()
}

func (s *FileStore) findIndex(name string) int {
	for i, rec := range s.cache {
		if rec.Name == name {
			return i
		}
	}

	return -1
}

func (s *FileStore) loadLocked() error {
	if s.loaded {
		return nil
	}
	s.loaded = true

	raw, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}

	key, err := s.deriveKey()
	if err != nil {
		return err
	}

	plain, err := utils.Decrypt(key, raw)
	if err != nil {
		return fmt.Errorf("decrypt secrets file (wrong TRYIT_SECRETS_PASSPHRASE?): %w", err)
	}

	var v fileVault
	if err := json.Unmarshal(plain, &v); err != nil {
		return fmt.Errorf("secrets file plaintext is not valid JSON: %w", err)
	}

	s.cache = v.Secrets
	return nil
}

func (s *FileStore) flushLocked() error {
	key, err := s.deriveKey()
	if err != nil {
		return err
	}

	plain, err := json.Marshal(fileVault{Secrets: s.cache})
	if err != nil {
		return err
	}

	blob, err := utils.Encrypt(key, plain)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}

	return os.WriteFile(s.path, blob, 0o600)
}

func (s *FileStore) deriveKey() ([32]byte, error) {
	saltFull := sha256.Sum256([]byte(config.GetFileSaltSeed()))
	raw, err := utils.DeriveKey(
		s.passphrase,
		saltFull[:16],
		config.GetScryptN(),
		config.GetScryptR(),
		config.GetScryptP(),
		config.GetFileKeySize(),
	)

	if err != nil {
		return [32]byte{}, err
	}

	var key [32]byte
	copy(key[:], raw)

	return key, nil
}

func cloneRecord(rec *secretType.StoredSecret) *secretType.StoredSecret {
	c := *rec
	return &c
}
