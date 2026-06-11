package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// NewToken returns a 256-bit random token as a hex string with the configured
// "tk_" prefix so it is visually distinct in logs/UI.
func NewToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return config.GetTokenPrefix() + hex.EncodeToString(b[:]), nil
}
