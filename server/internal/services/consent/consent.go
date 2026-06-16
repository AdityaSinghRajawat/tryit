package consent

import (
	"errors"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	consentType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/consent"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (s *ConsentService) IsConsentGranted(secret, host string) bool {
	if secret == "" || host == "" {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.findIndex(secret, host) >= 0
}

func (s *ConsentService) GrantConsent(secret, host string) *config.CustomError {
	secret = strings.TrimSpace(secret)
	host = strings.TrimSpace(host)
	if secret == "" || host == "" {
		return config.NewCustomError(
			errors.New("secret and host are required"),
			config.GetErrCodeInvalidRequest(),
		)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.findIndex(secret, host) >= 0 {
		return nil
	}

	s.records = append(s.records, consentType.Record{
		Secret:    secret,
		Host:      host,
		GrantedAt: utils.GetCurrTimeStamp(),
	})
	if err := s.flush(); err != nil {
		return config.NewCustomError(err, config.GetErrCodeInternal())
	}

	return nil
}
