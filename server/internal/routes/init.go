// Package routes is the composition root for the HTTP layer: it constructs
// services + handlers, then mounts them with the middleware chain
// (recover → cors → security). cmd/app.go only calls config.Init, invokes
// NewRoutes, and runs the http.Server.
package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	executeHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/execute"
	healthHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/health"
	pairHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/pair"
	"github.com/AdityaSinghRajawat/tryit/server/internal/middlewares"
	executeSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/execute"
	pairSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/pair"
	secretSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// NewRoutes builds every dependency in order (services → handlers) and
// returns a ready-to-serve http.Handler.
func NewRoutes(log *utils.Logger) (http.Handler, error) {
	// Services own their own storage / state internally.
	pairService, err := pairSvc.NewPairService()
	if err != nil {
		return nil, err
	}
	secretService := secretSvc.NewSecretService()

	// Outbound HTTP client — injected into the execute service. BaseURL is
	// empty because the user's RequestSpec carries the absolute URL.
	httpClient := utils.NewHttpClient("", config.GetExecTimeout())
	executeService := executeSvc.NewExecuteService(secretService, httpClient)

	// Handlers receive their service dependencies as concrete types.
	healthH := healthHandler.NewHealthHandler(pairService)
	pairH := pairHandler.NewPairHandler(pairService)
	executeH := executeHandler.NewExecuteHandler(executeService)

	// Router + middleware chain. Order matters: recover wraps everything,
	// cors handles preflight + ACAO, security enforces Host + Origin + bearer.
	r := chi.NewRouter()
	r.Use(middlewares.Recoverer(log))
	r.Use(middlewares.CORSMiddleware(pairService))
	r.Use(middlewares.SecurityMiddleware(pairService, config.GetHostHeader()))

	r.Get(config.GetRoutePathHealth(), healthH.Get)
	r.Post(config.GetRoutePathPair(), pairH.Post)
	r.Post(config.GetRoutePathExecute(), executeH.Post)

	return r, nil
}
