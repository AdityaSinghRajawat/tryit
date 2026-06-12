package config

type routeConsts struct {
	pathHealth  string
	pathPair    string
	pathExecute string
	pathParse   string
}

var routeI = &routeConsts{
	pathHealth:  "/health",
	pathPair:    "/pair",
	pathExecute: "/execute",
	pathParse:   "/parse",
}

func GetRoutePathHealth() string  { return routeI.pathHealth }
func GetRoutePathPair() string    { return routeI.pathPair }
func GetRoutePathExecute() string { return routeI.pathExecute }
func GetRoutePathParse() string   { return routeI.pathParse }
