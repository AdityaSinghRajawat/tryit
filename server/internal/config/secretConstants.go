package config

type secretConsts struct {
	envPrefix   string
	userSuffix  string
	passSuffix  string
}

var secretI = &secretConsts{
	envPrefix:  "TRYIT_SECRET_",
	userSuffix: "_USER",
	passSuffix: "_PASS",
}

func GetSecretEnvPrefix() string  { return secretI.envPrefix }
func GetSecretUserSuffix() string { return secretI.userSuffix }
func GetSecretPassSuffix() string { return secretI.passSuffix }
