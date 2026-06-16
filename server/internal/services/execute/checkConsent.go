package execute

import (
	"net/http"

	executeType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/execute"
	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
)

// checkConsent returns the first un-granted (storedSecretName, host) pair,
// or nil when every secret referenced by the request is allowed against the
// target host. spec.SecretRefs() reads only Auth.* placeholders — the only
// place the contract permits secrets.
func (s *ExecuteService) checkConsent(
	req *http.Request,
	spec specType.RequestSpec,
	mappedRefs map[string]string,
) *executeType.ConsentRequired {
	host := req.URL.Hostname()
	if host == "" {
		return nil
	}
	for _, placeholder := range spec.SecretRefs() {
		stored := placeholder
		if mappedRefs != nil {
			if v, ok := mappedRefs[placeholder]; ok && v != "" {
				stored = v
			}
		}
		if !s.ConsentService.IsConsentGranted(stored, host) {
			return &executeType.ConsentRequired{Secret: stored, Host: host}
		}
	}
	return nil
}
