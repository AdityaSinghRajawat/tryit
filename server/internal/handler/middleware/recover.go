package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/logger"
)

// Recover catches panics and writes a redacted 500. The panic value is logged
// (it may contain secrets in a bug) but only "internal" reaches the client.
func Recover(log *logger.Logger) func(http.Handler) http.Handler {
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
