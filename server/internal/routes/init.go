// Package routes is the HTTP composition root.
package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdityaSinghRajawat/tryit/server/api"
	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	consentHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/consent"
	executeHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/execute"
	healthHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/health"
	pairHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/pair"
	parseHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/parse"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai"
	"github.com/AdityaSinghRajawat/tryit/server/internal/middlewares"
	consentSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/consent"
	executeSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/execute"
	pairSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/pair"
	parseSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/parse"
	secretSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
	"github.com/AdityaSinghRajawat/tryit/server/internal/validations"
)

func NewRoutes() (http.Handler, error) {
	pairService, err := pairSvc.NewPairService()
	if err != nil {
		return nil, err
	}
	secretService := secretSvc.NewSecretService()
	consentService, cErr := consentSvc.NewConsentService(config.GetConsentFile())
	if cErr != nil {
		return nil, cErr
	}

	httpClient := utils.NewHttpClient("", config.GetExecTimeout())
	executeService := executeSvc.NewExecuteService(secretService, consentService, httpClient)

	validator, vErr := validations.NewSchemaValidator(api.Schema)
	if vErr != nil {
		return nil, vErr
	}
	aiProvider, aiErr := ai.NewAIProvider()
	if aiErr != nil {
		return nil, aiErr
	}
	cache := utils.NewCache()
	parseService := parseSvc.NewParseService(aiProvider, cache, validator)

	healthH := healthHandler.NewHealthHandler(pairService)
	pairH := pairHandler.NewPairHandler(pairService)
	executeH := executeHandler.NewExecuteHandler(executeService)
	parseH := parseHandler.NewParseHandler(parseService)
	consentH := consentHandler.NewConsentHandler(consentService)

	// Chain order matters: recover wraps everything, cors handles preflight,
	// security enforces Host + Origin + bearer.
	r := chi.NewRouter()
	r.Use(middlewares.Recoverer())
	r.Use(middlewares.CORSMiddleware(pairService))
	r.Use(middlewares.SecurityMiddleware(pairService, config.GetHostHeader()))

	r.Get(config.GetRoutePathHealth(), healthH.Get)
	r.Post(config.GetRoutePathPair(), pairH.Post)
	r.Post(config.GetRoutePathExecute(), executeH.Post)
	r.Post(config.GetRoutePathParse(), parseH.Post)
	r.Post(config.GetRoutePathConsent(), consentH.Post)

	return r, nil
}
