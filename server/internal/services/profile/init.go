// Package profile holds builtin + learned site profiles. Builtin ships in
// api/profiles.json; learned is opt-in and persisted at TRYIT_PROFILES_FILE
// (default ~/.tryit/profiles.json). Learned overrides builtin per host.
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	profileType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/profile"
)

type ProfileService struct {
	path    string
	mu      sync.RWMutex
	learned []profileType.SiteProfile
	builtin []profileType.SiteProfile
}

func NewProfileService(path string, builtinJSON []byte) (*ProfileService, error) {
	if path == "" {
		return nil, errors.New("profiles file path is empty (set TRYIT_PROFILES_FILE or $HOME)")
	}
	s := &ProfileService{path: path}
	if err := s.loadBuiltin(builtinJSON); err != nil {
		return nil, fmt.Errorf("load builtin profiles: %w", err)
	}
	if err := s.load(); err != nil {
		return nil, fmt.Errorf("load learned profiles: %w", err)
	}
	return s, nil
}

func (s *ProfileService) loadBuiltin(raw []byte) error {
	if len(raw) == 0 {
		return nil
	}
	var arr []profileType.SiteProfile
	if err := json.Unmarshal(raw, &arr); err != nil {
		return err
	}
	for i := range arr {
		arr[i].Source = profileType.SourceBuiltin
	}
	s.builtin = arr
	return nil
}
