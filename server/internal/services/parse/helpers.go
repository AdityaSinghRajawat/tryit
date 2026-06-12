package parse

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	parseType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/parse"
	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
)

func (s *ParseService) cacheKey(req parseType.Request) string {
	body := []byte(req.ScopedMarkdown)
	if len(body) == 0 && req.StructuredHint != nil {
		if b, err := json.Marshal(req.StructuredHint); err == nil {
			body = b
		}
	}
	inner := sha256.Sum256(body)
	outer := sha256.Sum256([]byte(req.PageURL + "\n" + hex.EncodeToString(inner[:])))
	return config.GetCacheKeyPrefix() + hex.EncodeToString(outer[:])
}

func (s *ParseService) getCachedResponse(ctx context.Context, key string) *parseType.Response {
	raw, err := s.Redis.GetKey(ctx, key)
	if err != nil || raw == "" {
		return nil
	}
	var r parseType.Response
	if json.Unmarshal([]byte(raw), &r) != nil {
		return nil
	}
	r.Source = parseType.SourceCache
	return &r
}

func (s *ParseService) saveCachedResponse(ctx context.Context, key string, resp *parseType.Response) {
	if resp == nil {
		return
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		return
	}
	_ = s.Redis.SetKey(ctx, key, string(raw), config.GetCacheTTL())
}

// detectSpec accepts either { "requestSpec": {...} } or a bare RequestSpec.
func (s *ParseService) detectSpec(hint any) *specType.RequestSpec {
	if hint == nil {
		return nil
	}
	b, err := json.Marshal(hint)
	if err != nil {
		return nil
	}

	var wrapped struct {
		RequestSpec *specType.RequestSpec `json:"requestSpec"`
	}
	if err := json.Unmarshal(b, &wrapped); err == nil && wrapped.RequestSpec != nil {
		if wrapped.RequestSpec.Validate() == nil {
			return wrapped.RequestSpec
		}
	}

	var bare specType.RequestSpec
	if err := json.Unmarshal(b, &bare); err == nil && bare.Validate() == nil {
		return &bare
	}
	return nil
}

func (s *ParseService) buildSystemPrompt() string {
	return config.GetParseSystemPromptTemplate() + string(s.Validator.Raw())
}

func (s *ParseService) buildUserMessage(req parseType.Request) string {
	return "PAGE_URL: " + req.PageURL + "\n" +
		"ENDPOINT_DOC:\n" + req.ScopedMarkdown + "\n" +
		"AUTHENTICATION_DOC (may be empty):\n" + req.AuthSectionMarkdown
}
