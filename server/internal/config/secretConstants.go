package config

type secretConsts struct {
	envPrefix  string
	userSuffix string
	passSuffix string

	envKindBearer string
	envKindUser   string
	envKindPass   string

	providerEnv      string
	providerKeychain string
	providerFile     string

	keychainServiceName string
	keychainIndexKey    string

	fileSaltSeed string
	scryptN      int
	scryptR      int
	scryptP      int
	fileKeySize  int
}

var secretI = &secretConsts{
	envPrefix:  "TRYIT_SECRET_",
	userSuffix: "_USER",
	passSuffix: "_PASS",

	envKindBearer: "bearer",
	envKindUser:   "user",
	envKindPass:   "pass",

	providerEnv:      "env",
	providerKeychain: "keychain",
	providerFile:     "file",

	keychainServiceName: "tryit-secrets",
	keychainIndexKey:    "__tryit_index__",

	fileSaltSeed: "tryit-secrets-salt-v1",
	scryptN:      32768,
	scryptR:      8,
	scryptP:      1,
	fileKeySize:  32,
}

func GetSecretEnvPrefix() string  { return secretI.envPrefix }
func GetSecretUserSuffix() string { return secretI.userSuffix }
func GetSecretPassSuffix() string { return secretI.passSuffix }

func GetEnvKindBearer() string { return secretI.envKindBearer }
func GetEnvKindUser() string   { return secretI.envKindUser }
func GetEnvKindPass() string   { return secretI.envKindPass }

func GetSecretsProviderEnv() string      { return secretI.providerEnv }
func GetSecretsProviderKeychain() string { return secretI.providerKeychain }
func GetSecretsProviderFile() string     { return secretI.providerFile }

func GetKeychainServiceName() string { return secretI.keychainServiceName }
func GetKeychainIndexKey() string    { return secretI.keychainIndexKey }

func GetFileSaltSeed() string { return secretI.fileSaltSeed }
func GetScryptN() int         { return secretI.scryptN }
func GetScryptR() int         { return secretI.scryptR }
func GetScryptP() int         { return secretI.scryptP }
func GetFileKeySize() int     { return secretI.fileKeySize }
