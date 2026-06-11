package config

type securityConsts struct {
	pathHealth            string
	pathPair              string
	headerOrigin          string
	headerHost            string
	extensionOriginPrefix string
}

var securityI = &securityConsts{
	pathHealth:            "/health",
	pathPair:              "/pair",
	headerOrigin:          "Origin",
	headerHost:            "Host",
	extensionOriginPrefix: "chrome-extension://",
}

func GetPathHealth() string            { return securityI.pathHealth }
func GetPathPair() string              { return securityI.pathPair }
func GetHeaderOrigin() string          { return securityI.headerOrigin }
func GetHeaderHost() string            { return securityI.headerHost }
func GetExtensionOriginPrefix() string { return securityI.extensionOriginPrefix }

// IsExemptPath returns true when the path bypasses the bearer + origin gate
// (only /health and /pair).
func IsExemptPath(p string) bool {
	return p == securityI.pathHealth || p == securityI.pathPair
}
