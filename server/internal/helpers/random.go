package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// Token returns a 256-bit random token as a 64-char hex string, prefixed
// "tk_" so it is visually distinct in logs/UI.
func Token() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "tk_" + hex.EncodeToString(b[:]), nil
}
