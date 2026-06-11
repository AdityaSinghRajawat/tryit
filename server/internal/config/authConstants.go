package config

import "regexp"

type authConsts struct {
	headerAuthorization string
	prefixBearer        string
	prefixBasic         string
	placeholderRe       *regexp.Regexp
}

var authI = &authConsts{
	headerAuthorization: "Authorization",
	prefixBearer:        "Bearer ",
	prefixBasic:         "Basic ",
	placeholderRe:       regexp.MustCompile(`\{\{secret:([A-Z0-9_]+)\}\}`),
}

func GetHeaderAuthorization() string            { return authI.headerAuthorization }
func GetPrefixBearer() string                   { return authI.prefixBearer }
func GetPrefixBasic() string                    { return authI.prefixBasic }
func GetSecretPlaceholderRegex() *regexp.Regexp { return authI.placeholderRe }
