package utils

import "strings"

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
