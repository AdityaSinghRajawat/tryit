package config

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// envConfig holds every runtime-configurable value (IMPL §10.1). The singleton
// is populated once by Init() at boot; all consumers read via getters.
type envConfig struct {
	port                int
	secretsBackend      string
	execTimeout         time.Duration
	execMaxResponseSize int64
	execInsecureTLS     bool
	logLevel            string
	pairFile            string
	homeDir             string
}

var envConfigI *envConfig

// Init populates the singleton from process env vars + defaults. Idempotent:
// a second call is a no-op. Any env var that fails to parse falls back to its
// default — getEnv*WithDefault helpers absorb the error.
func Init() error {
	if envConfigI != nil {
		return nil
	}
	homeDir, _ := os.UserHomeDir()
	envConfigI = &envConfig{
		port:                getEnvIntWithDefault("TRYIT_PORT", 8765),
		secretsBackend:      getEnvWithDefault("TRYIT_SECRETS_BACKEND", "env"),
		execTimeout:         getEnvDurationWithDefault("TRYIT_EXEC_TIMEOUT", 30*time.Second),
		execMaxResponseSize: getEnvInt64WithDefault("TRYIT_EXEC_MAX_RESPONSE_BYTES", 10*1024*1024),
		execInsecureTLS:     getEnvBoolWithDefault("TRYIT_EXEC_INSECURE_TLS", false),
		logLevel:            strings.ToLower(getEnvWithDefault("TRYIT_LOG_LEVEL", "info")),
		pairFile:            getEnvWithDefault("TRYIT_PAIR_FILE", defaultPairFile(homeDir)),
		homeDir:             homeDir,
	}
	return nil
}

func defaultPairFile(home string) string {
	if home == "" {
		return ""
	}
	return home + "/.tryit/pair.json"
}

// --- reusable env helpers --------------------------------------------------

// GetEnvByKey reads a runtime env var by key. Use this from services that
// need dynamically-named env vars (e.g. the secret service's per-NAME
// TRYIT_SECRET_<NAME> reads) — services should never touch os.Getenv
// directly. Returns the empty string if unset.
func GetEnvByKey(key string) string { return os.Getenv(key) }

func getEnvWithDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return v
	}
	return defaultValue
}

func getEnvInt64WithDefault(key string, defaultValue int64) int64 {
	if v, err := strconv.ParseInt(os.Getenv(key), 10, 64); err == nil {
		return v
	}
	return defaultValue
}

func getEnvBoolWithDefault(key string, defaultValue bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	if b, err := strconv.ParseBool(v); err == nil {
		return b
	}
	return defaultValue
}

func getEnvDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if d, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return d
	}
	return defaultValue
}

// --- getters --------------------------------------------------------------

func GetPort() int                  { return envConfigI.port }
func GetSecretsBackend() string     { return envConfigI.secretsBackend }
func GetExecTimeout() time.Duration { return envConfigI.execTimeout }
func GetExecMaxResponseSize() int64 { return envConfigI.execMaxResponseSize }
func GetExecInsecureTLS() bool      { return envConfigI.execInsecureTLS }
func GetLogLevel() string           { return envConfigI.logLevel }
func GetPairFile() string           { return envConfigI.pairFile }
func GetHomeDir() string            { return envConfigI.homeDir }

func GetListenAddr() string { return net.JoinHostPort("127.0.0.1", strconv.Itoa(envConfigI.port)) }
func GetHostHeader() string { return "127.0.0.1:" + strconv.Itoa(envConfigI.port) }
