// http.go is the HTTP utility hub: the injectable HttpClient (constructed
// once by routes, held by the service), request building from a RequestSpec,
// and the small header/body shaping helpers callers use to render previews.
// Nothing here is service-level business logic — services compose these.
package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
)

// HTTPResponse is the shaped outcome of an outbound HTTP call after body cap
// + truncation tracking. Services map it into wire DTOs.
type HTTPResponse struct {
	Status    int
	Headers   http.Header
	Body      []byte
	Truncated bool
}

// ----- HttpClient (injectable) -------------------------------------------

// HttpClient wraps *http.Client with config-driven timeouts, a TLS policy,
// and a redirect policy that strips Authorization on cross-host hops so a
// secret can never leak to a redirected target. Construct one in
// routes.NewRoutes and inject it into the services that need outbound HTTP.
type HttpClient struct {
	client  *http.Client
	BaseURL string // optional — Phase 2 named integrations set this; the dynamic execute path leaves it empty.
}

// NewHttpClient builds an HttpClient with the supplied total timeout and
// optional baseURL (empty when the caller will pass absolute URLs).
func NewHttpClient(baseURL string, timeout time.Duration) *HttpClient {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: timeout,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          50,
		IdleConnTimeout:       90 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: config.GetExecInsecureTLS(), // #nosec G402 — opt-in via env
		},
	}

	return &HttpClient{
		client: &http.Client{
			Transport:     tr,
			Timeout:       timeout,
			CheckRedirect: stripAuthOnCrossHost,
		},
		BaseURL: baseURL,
	}
}

// Do sends the prepared request and returns the shaped response. Body reads
// are capped at config.GetExecMaxResponseSize(); transport errors surface as
// errors and HTTP statuses are returned as data.
func (c *HttpClient) Do(req *http.Request) (*HTTPResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, truncated, err := readCappedBody(resp.Body, config.GetExecMaxResponseSize())
	if err != nil {
		return nil, err
	}
	return &HTTPResponse{
		Status:    resp.StatusCode,
		Headers:   resp.Header.Clone(),
		Body:      body,
		Truncated: truncated,
	}, nil
}

// stripAuthOnCrossHost is the CheckRedirect policy: cap redirects per config
// and drop the Authorization header on any cross-host hop so a secret can
// never leak to a redirected target.
func stripAuthOnCrossHost(req *http.Request, via []*http.Request) error {
	if len(via) >= config.GetMaxRedirects() {
		return http.ErrUseLastResponse
	}
	if len(via) == 0 {
		return nil
	}
	src := via[len(via)-1].URL
	if !strings.EqualFold(src.Hostname(), req.URL.Hostname()) {
		req.Header.Del(config.GetHeaderAuthorization())
	}
	return nil
}

func readCappedBody(r io.Reader, limit int64) ([]byte, bool, error) {
	if limit <= 0 {
		b, err := io.ReadAll(r)
		return b, false, err
	}
	b, err := io.ReadAll(io.LimitReader(r, limit))
	if err != nil {
		return nil, false, err
	}
	var probe [1]byte
	n, perr := r.Read(probe[:])
	if errors.Is(perr, io.EOF) || n == 0 {
		return b, false, nil
	}
	return b, true, nil
}

// ----- Request building --------------------------------------------------

// BuildHTTPRequest produces an *http.Request from a RequestSpec, without
// secret resolution. Auth-related fields keep their {{secret:NAME}}
// placeholders verbatim; the caller (execute service) stamps secrets on the
// returned request as a separate step so secret values stay localised.
func BuildHTTPRequest(ctx context.Context, s specType.RequestSpec) (*http.Request, error) {
	u, err := url.Parse(s.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseUrl: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("baseUrl must be absolute (scheme + host)")
	}

	path := s.Path
	for _, p := range s.PathParams {
		path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(p.Value))
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u.Path = strings.TrimRight(u.Path, "/") + path

	q := u.Query()
	for _, p := range s.Query {
		if len(p.Values) > 0 {
			for _, v := range p.Values {
				q.Add(p.Name, v)
			}
			continue
		}
		if p.Value != "" || p.Required {
			q.Add(p.Name, p.Value)
		}
	}
	u.RawQuery = q.Encode()

	body, contentType, err := buildHTTPBody(s.Body)
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > config.GetMaxRequestBodyBytes() {
		return nil, fmt.Errorf("request body exceeds %d bytes", config.GetMaxRequestBodyBytes())
	}

	req, err := http.NewRequestWithContext(
		ctx,
		strings.ToUpper(s.Method),
		u.String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for _, h := range s.Headers {
		req.Header.Set(h.Name, h.Value)
	}
	if s.Body.ContentType != "" {
		req.Header.Set("Content-Type", s.Body.ContentType)
	}
	return req, nil
}

func buildHTTPBody(b specType.BodySpec) ([]byte, string, error) {
	switch specType.Encoding(b.Encoding) {
	case specType.EncodingNone, "":
		return nil, "", nil
	case specType.EncodingJSON:
		if len(b.JSON) == 0 {
			return []byte("null"), "application/json", nil
		}
		if !json.Valid(b.JSON) {
			return nil, "", errors.New("body.json is not valid JSON")
		}
		return []byte(b.JSON), "application/json", nil
	case specType.EncodingForm:
		v := url.Values{}
		for _, p := range b.Form {
			if len(p.Values) > 0 {
				for _, x := range p.Values {
					v.Add(p.Name, x)
				}
				continue
			}
			if p.Value != "" || p.Required {
				v.Add(p.Name, p.Value)
			}
		}
		return []byte(v.Encode()), "application/x-www-form-urlencoded", nil
	case specType.EncodingRaw:
		return []byte(b.Raw), "", nil
	case specType.EncodingMultipart:
		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		for _, p := range b.Form {
			vals := p.Values
			if len(vals) == 0 {
				vals = []string{p.Value}
			}
			for _, x := range vals {
				if err := mw.WriteField(p.Name, x); err != nil {
					return nil, "", err
				}
			}
		}
		if err := mw.Close(); err != nil {
			return nil, "", err
		}
		return buf.Bytes(), mw.FormDataContentType(), nil
	default:
		return nil, "", fmt.Errorf("unknown body.encoding %q", b.Encoding)
	}
}

// ----- Header / body shaping helpers -------------------------------------

// FlattenHeaders collapses a multi-valued http.Header into a single-value
// map[string]string keyed by header name. Empty header values are dropped.
func FlattenHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) == 0 {
			continue
		}
		out[k] = v[0]
	}
	return out
}

// RequestBodyPreview returns up to 64 KiB of an *http.Request's body without
// consuming it. The request must have been built with a seekable body
// (http.NewRequestWithContext + bytes.NewReader sets GetBody automatically);
// otherwise an empty string is returned.
func RequestBodyPreview(req *http.Request) string {
	if req.Body == nil || req.GetBody == nil {
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

// StripHeaderName drops the "Name: " prefix from a "Name: Value" string and
// returns only the value half. Returns the input unchanged if no colon is
// present. Used for masked-preview Authorization rendering.
func StripHeaderName(headerLine string) string {
	for i := 0; i < len(headerLine); i++ {
		if headerLine[i] == ':' && i+2 <= len(headerLine) {
			return headerLine[i+2:]
		}
	}
	return headerLine
}
