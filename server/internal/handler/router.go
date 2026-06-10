package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdityaSinghRajawat/tryit/server/internal/handler/middleware"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/logger"
)

type Deps struct {
	Pair    middleware.Pairing
	Health  http.Handler
	Pairing http.Handler
	Execute http.Handler
	Logger  *logger.Logger
	Host    string // "127.0.0.1:<port>", for the Host-header check
}

// NewRouter wires routes + middleware in order:
//   1. recover (catches everything)
//   2. cors  (preflight + ACAO header)
//   3. security (Host + Origin + bearer; exempts /health and /pair)
//   4. dispatch
func NewRouter(d Deps) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recover(d.Logger))
	r.Use(middleware.CORS(d.Pair))
	r.Use(middleware.Security(d.Pair, d.Host))

	r.Method(http.MethodGet, "/health", d.Health)
	r.Method(http.MethodPost, "/pair", d.Pairing)
	r.Method(http.MethodPost, "/execute", d.Execute)
	return r
}
