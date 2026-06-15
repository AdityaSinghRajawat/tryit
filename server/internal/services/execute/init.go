// Package execute orchestrates build → consent-check → inject → send → mask
// for a user-supplied RequestSpec.
package execute

import (
	consentSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/consent"
	secretSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

type ExecuteService struct {
	SecretService  *secretSvc.SecretService
	ConsentService *consentSvc.ConsentService
	HttpClient     *utils.HttpClient
}

func NewExecuteService(
	secretService *secretSvc.SecretService,
	consentService *consentSvc.ConsentService,
	httpClient *utils.HttpClient,
) *ExecuteService {
	return &ExecuteService{
		SecretService:  secretService,
		ConsentService: consentService,
		HttpClient:     httpClient,
	}
}
