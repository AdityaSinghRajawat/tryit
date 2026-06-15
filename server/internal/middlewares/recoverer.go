// Package middlewares: flat HTTP middlewares as factory closures.
package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// Recoverer logs the panic value + stack and returns the standard envelope
// via HandleCustomError so the wire shape stays consistent with every other
// error path.
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
					utils.HandleCustomError(w, config.NewCustomError(
						errors.New("internal server error"),
						config.GetErrCodeInternal(),
					))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
