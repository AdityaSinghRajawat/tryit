// Package pair owns the pairing handshake (IMPL §8.1) and the file-backed
// store for the token + bound origin.
package pair

import (
	"errors"
	"sync"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

type pairFile struct {
	Token       string `json:"token"`
	BoundOrigin string `json:"boundOrigin"`
}

type PairService struct {
	path string
	mu   sync.RWMutex
	data pairFile
}

// NewPairService opens (or creates) the pair file and prints the token to
// stdout once on construction.
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
