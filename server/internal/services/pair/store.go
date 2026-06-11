package pair

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (s *PairService) load() (freshlyGenerated bool, err error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return false, err
	}

	b, err := os.ReadFile(s.path)
	if err == nil {
		var f pairFile
		if jerr := json.Unmarshal(b, &f); jerr != nil {
			return false, jerr
		}
		s.data = f
		if f.Token == "" {
			return s.regenerate()
		}
		return false, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	return s.regenerate()
}

func (s *PairService) regenerate() (bool, error) {
	tok, err := utils.NewToken()
	if err != nil {
		return false, err
	}
	s.data = pairFile{Token: tok}
	return true, s.flush()
}

func (s *PairService) flush() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o600)
}

func (s *PairService) Token() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.Token
}

func (s *PairService) BoundOrigin() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.BoundOrigin
}

func (s *PairService) SetBoundOrigin(origin string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.BoundOrigin = origin
	return s.flush()
}

func (s *PairService) Reset() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.regenerate()
	return err
}

func (s *PairService) announce(fresh bool) {
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "──────────────────────── tryit pairing ────────────────────────")
	switch {
	case s.BoundOrigin() != "":
		fmt.Fprintf(os.Stdout, " Already paired with: %s\n", s.BoundOrigin())
	case fresh:
		fmt.Fprintln(os.Stdout, " Fresh pairing token (paste into the extension panel):")
	default:
		fmt.Fprintln(os.Stdout, " Existing pairing token (paste into the extension panel):")
	}
	fmt.Fprintf(os.Stdout, "   %s\n", s.Token())
	fmt.Fprintln(os.Stdout, " Reset:  make reset-pairing")
	fmt.Fprintln(os.Stdout, "───────────────────────────────────────────────────────────────")
}
