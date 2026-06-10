package execute

import (
	"context"
	"net/http"
	"time"

	"github.com/AdityaSinghRajawat/tryit/server/internal/model"
)

// TargetClient is the port for the outbound HTTP integration. The concrete
// impl lives at integration/target and imports this package only to satisfy
// the port (hexagonal: services own their ports).
type TargetClient interface {
	Do(*http.Request) (*Response, error)
}

// Response is the raw outcome from the target call, before masking/preview
// assembly. Held in the service package so the integration is the one that
// imports a port type — not the other way around.
type Response struct {
	Status    int
	Headers   http.Header
	Body      []byte
	Truncated bool
}

// Outcome is what executeService returns to the handler. The handler maps it
// to the wire DTO. RequestPreview already has secrets masked.
type Outcome struct {
	Status         int
	Headers        map[string]string
	Body           string
	Truncated      bool
	DurationMs     int64
	RequestPreview RequestPreview
}

type RequestPreview struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body,omitempty"`
}

type Service struct {
	resolver SecretResolver
	target   TargetClient
}

func New(resolver SecretResolver, target TargetClient) *Service {
	return &Service{resolver: resolver, target: target}
}

// Execute orchestrates: validate → build → inject → send → mask. Consent and
// AI-cascade are out of scope for Phase 1 (D-P1-6).
func (s *Service) Execute(ctx context.Context, spec model.RequestSpec, refs map[string]string) (*Outcome, error) {
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	req, err := buildRequest(ctx, spec)
	if err != nil {
		return nil, err
	}

	// Capture a pre-injection snapshot for the masked preview. The injector
	// writes secret-bearing headers ONLY into req; our preview reuses the
	// masked return value.
	preHeaders := cloneHeaders(req.Header)
	preURL := req.URL.String()
	preMethod := req.Method
	preBody := requestBodyOrEmpty(req)

	maskedAuth, err := injectAuth(req, spec, refs, s.resolver)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	res, err := s.target.Do(req)
	if err != nil {
		return nil, err
	}
	durMs := time.Since(start).Milliseconds()

	preview := RequestPreview{
		Method:  preMethod,
		URL:     maskedURL(preURL, req.URL.String(), spec),
		Headers: preHeaders,
		Body:    preBody,
	}
	if maskedAuth != "" {
		// Replace the real Authorization (set by injector) with the masked
		// version in the preview.
		preview.Headers["Authorization"] = stripHeaderName(maskedAuth)
	}

	return &Outcome{
		Status:         res.Status,
		Headers:        flattenHeaders(res.Headers),
		Body:           string(res.Body),
		Truncated:      res.Truncated,
		DurationMs:     durMs,
		RequestPreview: preview,
	}, nil
}

// maskedURL returns the pre-injection URL — the apiKey-in-query case would
// otherwise leak the value, and other auth types don't change the URL at all.
// We surface the apiKey separately in the masked headers map.
func maskedURL(preURL, _ string, _ model.RequestSpec) string { return preURL }

func cloneHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) == 0 {
			continue
		}
		out[k] = v[0]
	}
	return out
}

func flattenHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) == 0 {
			continue
		}
		out[k] = v[0]
	}
	return out
}

func requestBodyOrEmpty(req *http.Request) string {
	if req.Body == nil {
		return ""
	}
	// We constructed the body via bytes.NewReader in buildRequest, so
	// GetBody returns a fresh reader without consuming the original.
	if req.GetBody == nil {
		return ""
	}
	body, err := req.GetBody()
	if err != nil {
		return ""
	}
	defer body.Close()
	const previewCap = 1 << 16
	buf := make([]byte, previewCap)
	n, _ := body.Read(buf)
	return string(buf[:n])
}

func stripHeaderName(masked string) string {
	// masked looks like "Authorization: Bearer ••••4242"; we want only the
	// value half ("Bearer ••••4242"). For unknown shapes return as-is.
	for i := 0; i < len(masked); i++ {
		if masked[i] == ':' && i+2 <= len(masked) {
			return masked[i+2:]
		}
	}
	return masked
}
