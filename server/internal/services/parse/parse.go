package parse

import (
	"context"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	parseType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/parse"
)

func (s *ParseService) Parse(
	ctx context.Context,
	req parseType.Request,
) (*parseType.Response, *config.CustomError) {
	key := s.cacheKey(req)

	if !req.Force {
		if hit := s.getCachedResponse(ctx, key); hit != nil {
			return hit, nil
		}
	}

	if spec := s.detectSpec(req.StructuredHint); spec != nil {
		resp := parseType.BuildResponse(*spec, parseType.SourceSpec)
		s.saveCachedResponse(ctx, key, resp)
		return resp, nil
	}

	spec, cerr := s.generateRequestSpec(ctx, req)
	if cerr != nil {
		return nil, cerr
	}
	resp := parseType.BuildResponse(*spec, parseType.SourceAI)
	s.saveCachedResponse(ctx, key, resp)
	return resp, nil
}
