// Package middlewares holds HTTP middlewares as flat files. Each exports a
// factory that closes over its dependencies and returns the standard
// `func(http.Handler) http.Handler` shape so chi can mount it.
package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// Recoverer catches panics, logs the redacted context, and returns a generic
// 500. The panic value is logged (it may contain secrets in a bug) but the
// response is always the standard envelope so no detail leaks.
func Recoverer(log *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rv := recover(); rv != nil {
					log.Error("panic in handler",
						"path", r.URL.Path,
						"value", rv,
						"stack", string(debug.Stack()),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write(
						[]byte(`{"error":{"code":"internal","message":"internal server error"}}`),
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
