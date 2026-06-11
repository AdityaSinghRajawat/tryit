package config

type routeConsts struct {
	pathHealth  string
	pathPair    string
	pathExecute string
}

var routeI = &routeConsts{
	pathHealth:  "/health",
	pathPair:    "/pair",
	pathExecute: "/execute",
}

func GetRoutePathHealth() string  { return routeI.pathHealth }
func GetRoutePathPair() string    { return routeI.pathPair }
func GetRoutePathExecute() string { return routeI.pathExecute }
