package consent

import consentSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/consent"

type ConsentHandler struct {
	ConsentService *consentSvc.ConsentService
}

func NewConsentHandler(consentService *consentSvc.ConsentService) *ConsentHandler {
	return &ConsentHandler{ConsentService: consentService}
}
