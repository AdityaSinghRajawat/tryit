package helpers

import "encoding/base64"

// BasicAuthValue returns the value half of a Basic-auth Authorization header:
// "Basic base64(user:pass)". No "Basic " prefix.
func BasicAuthValue(user, pass string) string {
	return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
}
