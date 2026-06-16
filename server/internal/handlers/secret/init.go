// Package secret serves the /secrets CRUD endpoints. Listing returns
// metadata only — values are never sent over the wire.
package secret

import (
	secretSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/secret"
)

type SecretHandler struct {
	SecretService *secretSvc.SecretService
}

func NewSecretHandler(secretService *secretSvc.SecretService) *SecretHandler {
	return &SecretHandler{SecretService: secretService}
}
