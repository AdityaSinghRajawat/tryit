package config

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type envConfig struct {
	// Server + execution
	port                int
	execTimeout         time.Duration
	execMaxResponseSize int64
	execInsecureTLS     bool
	logLevel            string

	// Secrets
	secretsBackend    string
	secretsFile       string
	secretsPassphrase string

	// Pairing
	pairFile string

	// AI provider — switchable via AI_PROVIDER; each provider has its own key.
	aiProvider      string
	openaiAPIKey    string
	anthropicAPIKey string
	geminiAPIKey    string
	ollamaHost      string
	ollamaModel     string

	// Cache (Redis-backed parse cache)
	cacheEnabled  bool
	cacheTTL      time.Duration
	redisAddr     string
	redisPassword string
	redisDB       int

	homeDir string
}

var envConfigI *envConfig

// Init populates the singleton from env vars + defaults. Idempotent.
func Init() error {
	if envConfigI != nil {
		return nil
	}
	homeDir, _ := os.UserHomeDir()
	envConfigI = &envConfig{
		port:                getEnvIntWithDefault("TRYIT_PORT", 8765),
		execTimeout:         getEnvDurationWithDefault("TRYIT_EXEC_TIMEOUT", 30*time.Second),
		execMaxResponseSize: getEnvInt64WithDefault("TRYIT_EXEC_MAX_RESPONSE_BYTES", 10*1024*1024),
		execInsecureTLS:     getEnvBoolWithDefault("TRYIT_EXEC_INSECURE_TLS", false),
		logLevel:            strings.ToLower(getEnvWithDefault("TRYIT_LOG_LEVEL", "info")),

		secretsBackend:    getEnvWithDefault("TRYIT_SECRETS_BACKEND", "env"),
		secretsFile:       getEnvWithDefault("TRYIT_SECRETS_FILE", defaultSecretsFile(homeDir)),
		secretsPassphrase: getEnvWithDefault("TRYIT_SECRETS_PASSPHRASE", ""),

		pairFile: getEnvWithDefault("TRYIT_PAIR_FILE", defaultPairFile(homeDir)),

		aiProvider:      strings.ToLower(getEnvWithDefault("AI_PROVIDER", "")),
		openaiAPIKey:    getEnvWithDefault("OPENAI_API_KEY", ""),
		anthropicAPIKey: getEnvWithDefault("ANTHROPIC_API_KEY", ""),
		geminiAPIKey:    getEnvWithDefault("GEMINI_API_KEY", ""),
		ollamaHost:      getEnvWithDefault("OLLAMA_HOST", aiI.defaultOllamaHost),
		ollamaModel:     getEnvWithDefault("OLLAMA_MODEL", aiI.defaultOllamaModel),

		cacheEnabled:  getEnvBoolWithDefault("TRYIT_CACHE_ENABLED", true),
		cacheTTL:      getEnvDurationWithDefault("TRYIT_CACHE_TTL", 24*time.Hour),
		redisAddr:     getEnvWithDefault("TRYIT_REDIS_ADDR", ""),
		redisPassword: getEnvWithDefault("TRYIT_REDIS_PASSWORD", ""),
		redisDB:       getEnvIntWithDefault("TRYIT_REDIS_DB", 0),

		homeDir: homeDir,
	}
	return nil
}

func defaultPairFile(home string) string {
	if home == "" {
		return ""
	}
	return home + "/.tryit/pair.json"
}

func defaultSecretsFile(home string) string {
	if home == "" {
		return ""
	}
	return home + "/.tryit/secrets.enc"
}

// GetEnvByKey is the dynamic-key escape hatch — services should never touch
// os.Getenv directly.
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

func GetPort() int                  { return envConfigI.port }
func GetExecTimeout() time.Duration { return envConfigI.execTimeout }
func GetExecMaxResponseSize() int64 { return envConfigI.execMaxResponseSize }
func GetExecInsecureTLS() bool      { return envConfigI.execInsecureTLS }
func GetLogLevel() string           { return envConfigI.logLevel }

func GetSecretsBackend() string    { return envConfigI.secretsBackend }
func GetSecretsFile() string       { return envConfigI.secretsFile }
func GetSecretsPassphrase() string { return envConfigI.secretsPassphrase }

func GetPairFile() string { return envConfigI.pairFile }

func GetAIProvider() string      { return envConfigI.aiProvider }
func GetOpenAIAPIKey() string    { return envConfigI.openaiAPIKey }
func GetAnthropicAPIKey() string { return envConfigI.anthropicAPIKey }
func GetGeminiAPIKey() string    { return envConfigI.geminiAPIKey }
func GetOllamaHost() string      { return envConfigI.ollamaHost }
func GetOllamaModel() string     { return envConfigI.ollamaModel }

func GetCacheEnabled() bool      { return envConfigI.cacheEnabled }
func GetCacheTTL() time.Duration { return envConfigI.cacheTTL }
func GetRedisAddr() string       { return envConfigI.redisAddr }
func GetRedisPassword() string   { return envConfigI.redisPassword }
func GetRedisDB() int            { return envConfigI.redisDB }

func GetHomeDir() string { return envConfigI.homeDir }

func GetListenAddr() string { return net.JoinHostPort("127.0.0.1", strconv.Itoa(envConfigI.port)) }
func GetHostHeader() string { return "127.0.0.1:" + strconv.Itoa(envConfigI.port) }
