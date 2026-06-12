package utils

import "encoding/base64"

// BasicAuthValue returns base64(user:pass), the value half of a Basic-auth header.
func BasicAuthValue(user, pass string) string {
	return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
}
