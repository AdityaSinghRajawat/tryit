package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// NewToken returns a "tk_<256-bit hex>" token.
func NewToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return config.GetTokenPrefix() + hex.EncodeToString(b[:]), nil
}
