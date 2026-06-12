// Package middlewares: flat HTTP middlewares as factory closures.
package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// Recoverer logs the panic value + stack and returns a generic 500 envelope.
func Recoverer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rv := recover(); rv != nil {
					utils.LogErrorWithStacktrace(
						r.Context(),
						fmt.Errorf("panic in handler: %v", rv),
						zap.String("path", r.URL.Path),
						zap.String("stack", string(debug.Stack())),
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
