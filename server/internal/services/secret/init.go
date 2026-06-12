// Package secret resolves NAMEd secrets. Env backend reads
// TRYIT_SECRET_<NAME> (plus _USER/_PASS for basic).
package secret

type SecretService struct{}

func NewSecretService() *SecretService {
	return &SecretService{}
}
