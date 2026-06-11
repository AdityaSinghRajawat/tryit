// Package pair owns both the pairing logic (IMPL §8.1) and the file-backed
// store for the token + bound origin (DECISION D-P1-3 — Phase 1 file at
// config.GetPairFile(); Phase 2 swaps to keychain behind the same surface).
//
// The store + business logic live in the same package because the surface is
// tiny and there's no DB layer to abstract.
package pair

import (
	"errors"
	"sync"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// pairFile is the on-disk JSON shape; unexported because nothing outside this
// package needs it.
type pairFile struct {
	Token       string `json:"token"`
	BoundOrigin string `json:"boundOrigin"`
}

// PairService holds the in-memory cache + file location and serves both the
// pairing-verify operation (Verify) and the read/write surface the
// middlewares need (Token + BoundOrigin).
type PairService struct {
	path string
	mu   sync.RWMutex
	data pairFile
}

// NewPairService opens (or creates) the pair file and prints the pairing
// token to stdout once on construction. Returns an error if the file path is
// unresolved.
func NewPairService() (*PairService, error) {
	path := config.GetPairFile()
	if path == "" {
		return nil, errors.New("pair store path is empty (set TRYIT_PAIR_FILE or $HOME)")
	}
	s := &PairService{path: path}
	fresh, err := s.load()
	if err != nil {
		return nil, err
	}
	s.announce(fresh)
	return s, nil
}
