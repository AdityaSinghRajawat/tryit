package pair

import (
	"crypto/subtle"
	"errors"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// Verify implements §8.1:
//   - If unbound: constant-time compare; on match, bind the calling Origin.
//   - If bound: require Origin == bound AND token match.
//
// Errors are classified inline as *config.CustomError so handlers can pass
// them straight to utils.HandleCustomError. Returns the (now-)bound origin
// on success.
func (s *PairService) Verify(token, origin string) (string, *config.CustomError) {
	if subtle.ConstantTimeCompare([]byte(strings.TrimSpace(token)), []byte(s.Token())) != 1 {
		return "", config.NewCustomError(
			errors.New("invalid token"),
			config.GetErrCodeUnauthorized(),
		)
	}

	if !strings.HasPrefix(origin, config.GetExtensionOriginPrefix()) {
		return "", config.NewCustomError(
			errors.New("origin not allowed"),
			config.GetErrCodeForbiddenOrigin(),
		)
	}

	bound := s.BoundOrigin()
	if bound == "" {
		if err := s.SetBoundOrigin(origin); err != nil {
			return "", config.NewCustomError(err, config.GetErrCodeInternal())
		}
		return origin, nil
	}

	if bound != origin {
		return "", config.NewCustomError(
			errors.New("origin conflicts with bound origin"),
			config.GetErrCodeForbiddenOrigin(),
		)
	}

	return bound, nil
}
