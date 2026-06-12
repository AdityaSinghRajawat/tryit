// Package middlewares: flat HTTP middlewares as factory closures.
package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// Recoverer logs the panic value + stack and returns a generic 500 envelope.
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
