// Package stores provides SecretStore implementations selected at boot.
package stores

import (
	"errors"
	"os"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
)

// EnvStore is read-only and backed by process env vars of the form
// TRYIT_SECRET_<NAME>, TRYIT_SECRET_<NAME>_USER, TRYIT_SECRET_<NAME>_PASS.
type EnvStore struct{}

func NewEnvStore() *EnvStore { return &EnvStore{} }

func (s *EnvStore) Get(name string) (*secretType.StoredSecret, error) {
	prefix := config.GetSecretEnvPrefix() + name

	if user := config.GetEnvByKey(prefix + config.GetSecretUserSuffix()); user != "" {
		pass := config.GetEnvByKey(prefix + config.GetSecretPassSuffix())
		return &secretType.StoredSecret{Name: name, Type: "basic", User: user, Pass: pass}, nil
	}

	if v := config.GetEnvByKey(prefix); v != "" {
		return &secretType.StoredSecret{Name: name, Type: "bearer", Value: v}, nil
	}

	return nil, nil
}

func (s *EnvStore) List() ([]*secretType.StoredSecret, error) {
	collected := map[string]*secretType.StoredSecret{}

	for _, kv := range os.Environ() {
		name, kind, value, ok := classifyEnvVar(kv)
		if !ok {
			continue
		}
		mergeEnvVar(collected, name, kind, value)
	}

	out := make([]*secretType.StoredSecret, 0, len(collected))
	for _, rec := range collected {
		out = append(out, rec)
	}

	return out, nil
}

func (s *EnvStore) Create(_ *secretType.StoredSecret) error {
	return errors.New("env secret store is read-only")
}

func (s *EnvStore) Delete(_ string) (bool, error) {
	return false, errors.New("env secret store is read-only")
}

// classifyEnvVar splits "KEY=VALUE" into (secretName, kind, value). kind is
// one of config.GetEnvKindBearer() / config.GetEnvKindUser() / config.GetEnvKindPass(). ok=false means the var
// is not a tryit secret.
func classifyEnvVar(kv string) (name, kind, value string, ok bool) {
	eq := strings.IndexByte(kv, '=')
	if eq < 0 {
		return "", "", "", false
	}

	prefix := config.GetSecretEnvPrefix()
	key := kv[:eq]
	if !strings.HasPrefix(key, prefix) {
		return "", "", "", false
	}

	rest := key[len(prefix):]
	value = kv[eq+1:]

	switch {
	case strings.HasSuffix(rest, config.GetSecretUserSuffix()):
		return strings.TrimSuffix(
			rest,
			config.GetSecretUserSuffix(),
		), config.GetEnvKindUser(), value, true
	case strings.HasSuffix(rest, config.GetSecretPassSuffix()):
		return strings.TrimSuffix(
			rest,
			config.GetSecretPassSuffix(),
		), config.GetEnvKindPass(), value, true
	default:
		return rest, config.GetEnvKindBearer(), value, true
	}
}

func mergeEnvVar(into map[string]*secretType.StoredSecret, name, kind, value string) {
	entry := into[name]

	switch kind {
	case config.GetEnvKindUser():
		if entry == nil {
			entry = &secretType.StoredSecret{Name: name}
			into[name] = entry
		}
		entry.Type = "basic"
		entry.User = value
	case config.GetEnvKindPass():
		if entry == nil {
			entry = &secretType.StoredSecret{Name: name}
			into[name] = entry
		}
		entry.Type = "basic"
		entry.Pass = value
	case config.GetEnvKindBearer():
		if entry != nil {
			return
		}
		into[name] = &secretType.StoredSecret{Name: name, Type: "bearer"}
	}
}
