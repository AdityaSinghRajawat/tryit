package middlewares

import (
	"net/http"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// CORSMiddleware answers preflight (OPTIONS) and attaches
// Access-Control-Allow-Origin for the bound extension origin. Never "*".
// First-time pairing (/pair) preflight is permitted for any chrome-extension://
// origin so the panel can submit its token before binding.
func CORSMiddleware(pair PairReader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get(config.GetHeaderOrigin())
			allow := ""
			switch {
			case pair.BoundOrigin() != "" && origin == pair.BoundOrigin():
				allow = origin
			case r.URL.Path == config.GetPathPair() && isExtensionOrigin(origin):
				allow = origin
			case r.URL.Path == config.GetPathHealth() && isExtensionOrigin(origin):
				allow = origin
			}
			if allow != "" {
				w.Header().Set("Access-Control-Allow-Origin", allow)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "600")
			}
			if r.Method == http.MethodOptions {
				if allow == "" {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isExtensionOrigin(origin string) bool {
	return strings.HasPrefix(origin, config.GetExtensionOriginPrefix())
}
