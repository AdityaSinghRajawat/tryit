package profile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	profileType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/profile"
)

func (s *ProfileService) load() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(b) == 0 {
		return nil
	}
	var arr []profileType.SiteProfile
	if err := json.Unmarshal(b, &arr); err != nil {
		return err
	}
	for i := range arr {
		arr[i].Source = profileType.SourceLearned
	}
	s.learned = arr
	return nil
}

func (s *ProfileService) flush() error {
	b, err := json.MarshalIndent(s.learned, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o600)
}

func (s *ProfileService) findLearnedIndex(host string) int {
	for i, p := range s.learned {
		if strings.EqualFold(p.Host, host) {
			return i
		}
	}
	return -1
}

func (s *ProfileService) findBuiltin(host string) *profileType.SiteProfile {
	for i := range s.builtin {
		if strings.EqualFold(s.builtin[i].Host, host) {
			return &s.builtin[i]
		}
	}
	return nil
}
