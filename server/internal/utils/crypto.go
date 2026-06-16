package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

const NonceLen = 24

// Encrypt seals plaintext under key with a fresh random nonce and returns
// the layout: nonce(24 bytes) || ciphertext. Decrypt reverses it.
func Encrypt(key [32]byte, plaintext []byte) ([]byte, error) {
	var nonce [NonceLen]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}
	ct := secretbox.Seal(nil, plaintext, &nonce, &key)
	out := make([]byte, 0, NonceLen+len(ct))
	out = append(out, nonce[:]...)
	out = append(out, ct...)
	return out, nil
}

// Decrypt opens a blob produced by Encrypt. Returns an error on tamper or a
// wrong key.
func Decrypt(key [32]byte, blob []byte) ([]byte, error) {
	if len(blob) < NonceLen+secretbox.Overhead {
		return nil, errors.New("encrypted blob too short")
	}
	var nonce [NonceLen]byte
	copy(nonce[:], blob[:NonceLen])
	plain, ok := secretbox.Open(nil, blob[NonceLen:], &nonce, &key)
	if !ok {
		return nil, errors.New("decrypt failed (wrong key or tampered ciphertext)")
	}
	return plain, nil
}

// DeriveKey runs scrypt with caller-supplied cost params and returns a
// keyLen-byte key. Cost params and salt are the caller's choice so this
// helper stays a pure primitive (no config coupling).
func DeriveKey(passphrase string, salt []byte, n, r, p, keyLen int) ([]byte, error) {
	return scrypt.Key([]byte(passphrase), salt, n, r, p, keyLen)
}

// NewToken returns a "tk_<256-bit hex>" token.
func NewToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return config.GetTokenPrefix() + hex.EncodeToString(b[:]), nil
}
