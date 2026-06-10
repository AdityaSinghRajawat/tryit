// Package helpers holds tiny pure utilities. Nothing here depends on services
// or handlers; nothing here imports a Secret value beyond its already-revealed
// string (callers are responsible).
package helpers

import "strings"

// Mask returns "••••XXXX" using the last 4 chars of v. Short values become
// "•••" (all redacted). Returns empty for empty.
func Mask(v string) string {
	if v == "" {
		return ""
	}
	if len(v) <= 4 {
		return "••••"
	}
	return "••••" + v[len(v)-4:]
}

// MaskBearer masks the credential half of "Bearer X" / "Basic X" / etc. while
// keeping the prefix readable, for header previews.
func MaskBearer(headerValue string) string {
	parts := strings.SplitN(headerValue, " ", 2)
	if len(parts) == 2 {
		return parts[0] + " " + Mask(parts[1])
	}
	return Mask(headerValue)
}
