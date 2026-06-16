package stores

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
	"github.com/zalando/go-keyring"
)

// KeychainStore backs SecretStore with the OS keyring via go-keyring.
// go-keyring has no enumerate API, so a sidecar index entry keyed by
// config.GetKeychainIndexKey() tracks the set of names.
type KeychainStore struct {
	serviceName string
	indexKey    string
}

func NewKeychainStore() *KeychainStore {
	return &KeychainStore{
		serviceName: config.GetKeychainServiceName(),
		indexKey:    config.GetKeychainIndexKey(),
	}
}

func (s *KeychainStore) Get(name string) (*secretType.StoredSecret, error) {
	raw, err := keyring.Get(s.serviceName, name)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var rec secretType.StoredSecret
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		return nil, err
	}

	rec.Name = name
	return &rec, nil
}

func (s *KeychainStore) List() ([]*secretType.StoredSecret, error) {
	names, err := s.readIndex()
	if err != nil {
		return nil, err
	}

	out := make([]*secretType.StoredSecret, 0, len(names))
	for _, n := range names {
		rec, err := s.Get(n)
		if err != nil {
			return nil, err
		}
		if rec != nil {
			out = append(out, rec)
		}
	}

	return out, nil
}

func (s *KeychainStore) Create(rec *secretType.StoredSecret) error {
	payload, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	if err := keyring.Set(s.serviceName, rec.Name, string(payload)); err != nil {
		return err
	}

	return s.addToIndex(rec.Name)
}

func (s *KeychainStore) Delete(name string) (bool, error) {
	err := keyring.Delete(s.serviceName, name)
	if errors.Is(err, keyring.ErrNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, s.removeFromIndex(name)
}

func (s *KeychainStore) addToIndex(name string) error {
	names, err := s.readIndex()
	if err != nil {
		return err
	}

	for _, n := range names {
		if n == name {
			return nil
		}
	}

	return s.writeIndex(append(names, name))
}

func (s *KeychainStore) removeFromIndex(name string) error {
	names, err := s.readIndex()
	if err != nil {
		return err
	}

	out := names[:0]
	for _, n := range names {
		if n != name {
			out = append(out, n)
		}
	}

	return s.writeIndex(out)
}

func (s *KeychainStore) readIndex() ([]string, error) {
	raw, err := keyring.Get(s.serviceName, s.indexKey)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	out := parts[:0]
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}

	return out, nil
}

func (s *KeychainStore) writeIndex(names []string) error {
	return keyring.Set(s.serviceName, s.indexKey, strings.Join(names, ","))
}
