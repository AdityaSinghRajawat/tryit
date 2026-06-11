// Package execute orchestrates the build → inject → send → mask flow for an
// outbound request driven by a user-supplied RequestSpec. HTTP plumbing
// (request building, sending, header/body shaping) lives in utils;
// secret-aware logic (auth resolution + stamping) lives here.
//
// Phase 1 skips the consent gate (D-P1-6); Phase 2 adds it.
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
