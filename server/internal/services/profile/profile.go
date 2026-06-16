package profile

import (
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	profileType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/profile"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// Lookup returns the profile for host, preferring learned over builtin.
// Returns nil when no profile matches.
func (s *ProfileService) Lookup(host string) *profileType.SiteProfile {
	host = strings.TrimSpace(host)
	if host == "" {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i := s.findLearnedIndex(host); i >= 0 {
		p := s.learned[i]
		return &p
	}
	if p := s.findBuiltin(host); p != nil {
		copy := *p
		return &copy
	}
	return nil
}

// Learn persists p as a learned profile. Replaces an existing learned entry
// for the same host; never modifies builtin.
func (s *ProfileService) Learn(p *profileType.SiteProfile) *config.CustomError {
	if err := p.Validate(); err != nil {
		return config.NewCustomError(err, config.GetErrCodeInvalidRequest())
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p.Source = profileType.SourceLearned
	p.UpdatedAt = utils.GetCurrTimeStamp()
	if i := s.findLearnedIndex(p.Host); i >= 0 {
		s.learned[i] = *p
	} else {
		s.learned = append(s.learned, *p)
	}
	if err := s.flush(); err != nil {
		return config.NewCustomError(err, config.GetErrCodeInternal())
	}
	return nil
}

// List returns the union of learned + builtin (learned overrides builtin for
// the same host). Stable iteration order: learned first, then builtin minus
// any host already covered by a learned entry.
func (s *ProfileService) List() []profileType.SiteProfile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]profileType.SiteProfile, 0, len(s.learned)+len(s.builtin))
	seen := make(map[string]struct{}, len(s.learned))
	for _, p := range s.learned {
		out = append(out, p)
		seen[strings.ToLower(p.Host)] = struct{}{}
	}
	for _, p := range s.builtin {
		if _, dup := seen[strings.ToLower(p.Host)]; dup {
			continue
		}
		out = append(out, p)
	}
	return out
}
