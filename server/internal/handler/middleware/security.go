// Package middleware implements the gate (§8.2):
//   - Host header anti-DNS-rebinding check
//   - Origin allowlist (bound at pair time, §8.1)
//   - Constant-time bearer token check
//   - 409 not_paired before a token is bound
//
// /health is fully exempt. /pair is exempt from token+origin (token lives in
// the body); the pair service does its own constant-time compare.
package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/dto"
)

// Pairing is the read side of pair.Store that middleware needs.
type Pairing interface {
	// Token returns the stored pairing token (always present from server start).
	Token() string
	// BoundOrigin returns the chrome-extension://<id> origin paired against,
	// or "" if the server has never been paired.
	BoundOrigin() string
}

// Exempt paths bypass token + origin checks.
var exempt = map[string]bool{
	"/health": true,
	"/pair":   true,
}

func Security(p Pairing, hostHeader string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS preflight is allowed through; cors middleware writes headers.
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			// Host check (anti-DNS-rebinding) — every request.
			if !hostMatches(r.Host, hostHeader) {
				writeErr(w, http.StatusForbidden, dto.ErrForbiddenOrigin, "host header rejected")
				return
			}

			if exempt[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			bound := p.BoundOrigin()
			if bound == "" {
				writeErr(w, http.StatusConflict, dto.ErrNotPaired, "server is not paired")
				return
			}
			if r.Header.Get("Origin") != bound {
				writeErr(w, http.StatusForbidden, dto.ErrForbiddenOrigin, "origin not allowed")
				return
			}

			tok := bearerToken(r.Header.Get("Authorization"))
			if subtle.ConstantTimeCompare([]byte(tok), []byte(p.Token())) != 1 {
				writeErr(w, http.StatusUnauthorized, dto.ErrUnauthorized, "invalid token")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// hostMatches accepts the configured "127.0.0.1:port" Host header. Browsers
// may also use "localhost:port" on some platforms — for the loopback bind we
// accept that variant too. Anything else is rejected (DNS rebinding defence).
func hostMatches(got, want string) bool {
	if got == want {
		return true
	}
	// Accept localhost: prefix as well — same loopback.
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

func writeErr(w http.ResponseWriter, status int, code dto.ErrorCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body := `{"error":{"code":"` + string(code) + `","message":"` + msg + `"}}`
	_, _ = w.Write([]byte(body))
}
