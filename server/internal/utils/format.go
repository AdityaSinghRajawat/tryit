package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// BasicAuthValue returns base64(user:pass), the value half of a Basic-auth header.
func BasicAuthValue(user, pass string) string {
	return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
}

// Mask returns "••••XXXX" with the last 4 chars of v, or "••••" if shorter.
func Mask(v string) string {
	if v == "" {
		return ""
	}
	if len(v) <= 4 {
		return "••••"
	}
	return "••••" + v[len(v)-4:]
}

// MaskBearer masks the credential half of "Bearer X" / "Basic X" / etc.
func MaskBearer(headerValue string) string {
	parts := strings.SplitN(headerValue, " ", 2)
	if len(parts) == 2 {
		return parts[0] + " " + Mask(parts[1])
	}
	return Mask(headerValue)
}

func PrettyJSON(raw json.RawMessage) string {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(raw)
	}
	return string(b)
}

func PlaceholderToEnv(ref string) (envName, secretName string) {
	re := config.GetSecretPlaceholderRegex()
	m := re.FindStringSubmatch(ref)
	if len(m) < 2 {
		return "", ""
	}
	return config.GetSecretEnvPrefix() + m[1], m[1]
}
