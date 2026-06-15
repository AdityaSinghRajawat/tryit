package middlewares

import (
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// PairReader is the narrow surface the cors + security middlewares need.
type PairReader interface {
	Token() string
	BoundOrigin() string
}

// SecurityMiddleware enforces IMPL §8.2: Host anti-DNS-rebinding, Origin must
// equal bound, constant-time bearer, 409 not_paired before binding. /health
// is fully exempt; /pair is exempt from token+origin (body token is the auth).
func SecurityMiddleware(pairReader PairReader, hostHeader string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			if !hostMatches(r.Host, hostHeader) {
				utils.HandleCustomError(w, config.NewCustomError(
					errors.New("host header rejected"),
					config.GetErrCodeForbiddenOrigin(),
				))
				return
			}
			if config.IsExemptPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			bound := pairReader.BoundOrigin()
			if bound == "" {
				utils.HandleCustomError(w, config.NewCustomError(
					errors.New("server is not paired"),
					config.GetErrCodeNotPaired(),
				))
				return
			}
			if r.Header.Get(config.GetHeaderOrigin()) != bound {
				utils.HandleCustomError(w, config.NewCustomError(
					errors.New("origin not allowed"),
					config.GetErrCodeForbiddenOrigin(),
				))
				return
			}
			tok := bearerToken(r.Header.Get(config.GetHeaderAuthorization()))
			if subtle.ConstantTimeCompare([]byte(tok), []byte(pairReader.Token())) != 1 {
				utils.HandleCustomError(w, config.NewCustomError(
					errors.New("invalid token"),
					config.GetErrCodeUnauthorized(),
				))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func hostMatches(got, want string) bool {
	if got == want {
		return true
	}
	if strings.HasPrefix(want, "127.0.0.1:") {
		_, port, ok := strings.Cut(want, ":")
		if ok && got == "localhost:"+port {
			return true
		}
	}
	return false
}

func bearerToken(h string) string {
	const p = "Bearer "
	if !strings.HasPrefix(h, p) {
		return ""
	}
	return strings.TrimSpace(h[len(p):])
}
