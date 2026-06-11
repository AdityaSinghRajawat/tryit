// Package secret resolves NAMEd secrets for the execute service. Phase 1
// backend reads TRYIT_SECRET_<NAME> env vars (plus _USER/_PASS for basic).
// Phase 2 swaps in keychain + encrypted-file backends behind the same shape.
package secret

type SecretService struct{}

func NewSecretService() *SecretService {
	return &SecretService{}
}
