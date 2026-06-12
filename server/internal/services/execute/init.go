// Package execute orchestrates build → inject → send → mask for a
// user-supplied RequestSpec. HTTP plumbing is in utils; secret-aware logic
// lives here.
package execute

import (
	secretSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

type ExecuteService struct {
	Resolver   *secretSvc.SecretService
	HttpClient *utils.HttpClient
}

func NewExecuteService(
	resolver *secretSvc.SecretService,
	httpClient *utils.HttpClient,
) *ExecuteService {
	return &ExecuteService{Resolver: resolver, HttpClient: httpClient}
}
