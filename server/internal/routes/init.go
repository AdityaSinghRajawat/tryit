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
	profileHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/profile"
	secretHandler "github.com/AdityaSinghRajawat/tryit/server/internal/handlers/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai"
	"github.com/AdityaSinghRajawat/tryit/server/internal/middlewares"
	consentSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/consent"
	executeSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/execute"
	pairSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/pair"
	parseSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/parse"
	profileSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/profile"
	secretSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
	"github.com/AdityaSinghRajawat/tryit/server/internal/validations"
)

func NewRoutes() (http.Handler, error) {
	pairService, err := pairSvc.NewPairService(config.GetPairFile())
	if err != nil {
		return nil, err
	}
	secretStore, sErr := secretSvc.NewSecretStore()
	if sErr != nil {
		return nil, sErr
	}
	secretService := secretSvc.NewSecretService(secretStore)
	consentService, cErr := consentSvc.NewConsentService(config.GetConsentFile())
	if cErr != nil {
		return nil, cErr
	}
	profileService, pErr := profileSvc.NewProfileService(
		config.GetProfilesFile(),
		api.BuiltinProfiles,
	)
	if pErr != nil {
		return nil, pErr
	}

	httpClient := utils.NewHttpClient("", config.GetExecTimeout())
	executeService := executeSvc.NewExecuteService(secretService, consentService, httpClient)

	schemaValidator, vErr := validations.NewSchemaValidator(api.Schema)
	if vErr != nil {
		return nil, vErr
	}
	aiProvider, aiErr := ai.NewAIProvider()
	if aiErr != nil {
		return nil, aiErr
	}
	cache := utils.NewCache(
		config.GetCacheLRUCapacity(),
		config.GetCacheTTL(),
		config.GetCacheDiskDir(),
		config.GetCacheEnabled(),
	)
	parseService := parseSvc.NewParseService(aiProvider, cache, schemaValidator, profileService)

	healthH := healthHandler.NewHealthHandler(pairService)
	pairH := pairHandler.NewPairHandler(pairService)
	executeH := executeHandler.NewExecuteHandler(executeService)
	parseH := parseHandler.NewParseHandler(parseService)
	consentH := consentHandler.NewConsentHandler(consentService)
	profileH := profileHandler.NewProfileHandler(profileService)
	secretH := secretHandler.NewSecretHandler(secretService)

	// Chain order matters: recover wraps everything, cors handles preflight,
	// security enforces Host + Origin + bearer.
	r := chi.NewRouter()
	r.Use(middlewares.Recoverer())
	r.Use(middlewares.CORSMiddleware(pairService))
	r.Use(middlewares.SecurityMiddleware(pairService, config.GetHostHeader()))

	r.Get(config.GetRoutePathHealth(), healthH.CheckHealth)
	r.Post(config.GetRoutePathPair(), pairH.CreatePair)
	r.Post(config.GetRoutePathExecute(), executeH.ExecuteCommand)
	r.Post(config.GetRoutePathParse(), parseH.ParseCommand)
	r.Post(config.GetRoutePathConsent(), consentH.CreateConsent)
	r.Get(config.GetRoutePathProfiles(), profileH.ListProfiles)
	r.Post(config.GetRoutePathProfiles(), profileH.CreateProfile)
	r.Get(config.GetRoutePathSecrets(), secretH.ListSecrets)
	r.Post(config.GetRoutePathSecrets(), secretH.CreateSecret)
	r.Delete(config.GetRoutePathSecretsByName(), secretH.DeleteSecret)

	return r, nil
}
