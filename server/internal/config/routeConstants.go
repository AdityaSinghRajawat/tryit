package config

type routeConsts struct {
	pathHealth   string
	pathPair     string
	pathExecute  string
	pathParse    string
	pathConsent  string
	pathProfiles string
}

var routeI = &routeConsts{
	pathHealth:   "/health",
	pathPair:     "/pair",
	pathExecute:  "/execute",
	pathParse:    "/parse",
	pathConsent:  "/consent",
	pathProfiles: "/profiles",
}

func GetRoutePathHealth() string   { return routeI.pathHealth }
func GetRoutePathPair() string     { return routeI.pathPair }
func GetRoutePathExecute() string  { return routeI.pathExecute }
func GetRoutePathParse() string    { return routeI.pathParse }
func GetRoutePathConsent() string  { return routeI.pathConsent }
func GetRoutePathProfiles() string { return routeI.pathProfiles }
