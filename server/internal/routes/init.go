// Package routes is the HTTP composition root.
package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"

	"github.com/AdityaSinghRajawat/tryit/server/api"
	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	executeHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/execute"
	healthHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/health"
	pairHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/pair"
	parseHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/parse"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai"
	"github.com/AdityaSinghRajawat/tryit/server/internal/middlewares"
	executeSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/execute"
	pairSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/pair"
	parseSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/parse"
	secretSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
	"github.com/AdityaSinghRajawat/tryit/server/internal/validations"
)

func NewRoutes(log *utils.Logger, redisClient *redis.Client) (http.Handler, error) {
	pairService, err := pairSvc.NewPairService()
	if err != nil {
		return nil, err
	}
	secretService := secretSvc.NewSecretService()

	httpClient := utils.NewHttpClient("", config.GetExecTimeout())
	executeService := executeSvc.NewExecuteService(secretService, httpClient)

	validator, vErr := validations.NewSchemaValidator(api.Schema)
	if vErr != nil {
		return nil, vErr
	}
	aiProvider, aiErr := ai.NewAIProvider()
	if aiErr != nil {
		return nil, aiErr
	}
	redisUtil := utils.NewRedisUtilManager(redisClient)
	parseService := parseSvc.NewParseService(aiProvider, redisUtil, validator)

	healthH := healthHandler.NewHealthHandler(pairService)
	pairH := pairHandler.NewPairHandler(pairService)
	executeH := executeHandler.NewExecuteHandler(executeService)
	parseH := parseHandler.NewParseHandler(parseService)

	// Chain order matters: recover wraps everything, cors handles preflight,
	// security enforces Host + Origin + bearer.
	r := chi.NewRouter()
	r.Use(middlewares.Recoverer(log))
	r.Use(middlewares.CORSMiddleware(pairService))
	r.Use(middlewares.SecurityMiddleware(pairService, config.GetHostHeader()))

	r.Get(config.GetRoutePathHealth(), healthH.Get)
	r.Post(config.GetRoutePathPair(), pairH.Post)
	r.Post(config.GetRoutePathExecute(), executeH.Post)
	r.Post(config.GetRoutePathParse(), parseH.Post)

	return r, nil
}
