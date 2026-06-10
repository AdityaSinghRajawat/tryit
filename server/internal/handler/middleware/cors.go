package middleware

import "net/http"

// CORS answers preflight (OPTIONS) and attaches Access-Control-Allow-Origin
// for the bound extension origin. Never "*". For routes hit before pairing,
// we echo the request Origin only if /pair (otherwise the Security middleware
// will reject anyway with 409 not_paired).
func CORS(p Pairing) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allow := ""
			switch {
			case p.BoundOrigin() != "" && origin == p.BoundOrigin():
				allow = origin
			case r.URL.Path == "/pair" && isExtensionOrigin(origin):
				// First-time pairing: the origin isn't bound yet. Allow
				// chrome-extension:// preflight to reach /pair so the panel
				// can submit the token.
				allow = origin
			case r.URL.Path == "/health":
				// /health is safe to expose; useful for the panel's bootstrap.
				if isExtensionOrigin(origin) {
					allow = origin
				}
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
	return len(origin) >= len("chrome-extension://") &&
		origin[:len("chrome-extension://")] == "chrome-extension://"
}
