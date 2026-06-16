package parse

import (
	"context"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	parseType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/parse"
)

func (s *ParseService) ParseCommand(
	ctx context.Context,
	req parseType.Request,
) (*parseType.Response, *config.CustomError) {
	key := s.cacheKey(req)

	if !req.Force {
		if hit := s.getCachedResponse(key); hit != nil {
			return hit, nil
		}
	}

	if spec := s.detectSpec(req.StructuredHint); spec != nil {
		resp := parseType.BuildResponse(*spec, parseType.SourceSpec)
		s.saveCachedResponse(key, resp)
		return resp, nil
	}

	spec, cerr := s.generateRequestSpec(ctx, req)
	if cerr != nil {
		return nil, cerr
	}

	source := parseType.SourceAI
	if prof := s.ProfileService.LookupProfile(hostFromPageURL(req.PageURL)); prof != nil {
		applyProfile(spec, prof)
		source = parseType.SourceProfile
	}

	resp := parseType.BuildResponse(*spec, source)
	s.saveCachedResponse(key, resp)

	return resp, nil
}
