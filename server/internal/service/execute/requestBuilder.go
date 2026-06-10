// Package execute builds → consent-checks → resolves → injects → sends →
// masks (IMPL §8.4, §9.1). requestBuilder is the pure spec→*http.Request
// transformation. authInjector is the secret-bearing half — separated so
// secret values exist in process memory only at the latest possible moment.
package execute

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/model"
)

const maxRequestBodyBytes = 5 * 1024 * 1024

// buildRequest produces an *http.Request from a RequestSpec, without secret
// resolution. Auth-related fields keep their {{secret:NAME}} placeholders
// verbatim; the injector replaces them later.
func buildRequest(ctx context.Context, spec model.RequestSpec) (*http.Request, error) {
	u, err := url.Parse(spec.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseUrl: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("baseUrl must be absolute (scheme + host)")
	}

	path := spec.Path
	for _, p := range spec.PathParams {
		path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(p.Value))
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u.Path = strings.TrimRight(u.Path, "/") + path

	q := u.Query()
	for _, p := range spec.Query {
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

	body, contentType, err := buildBody(spec.Body)
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > maxRequestBodyBytes {
		return nil, fmt.Errorf("request body exceeds %d bytes", maxRequestBodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(spec.Method), u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for _, h := range spec.Headers {
		req.Header.Set(h.Name, h.Value)
	}
	if spec.Body.ContentType != "" {
		req.Header.Set("Content-Type", spec.Body.ContentType)
	}
	return req, nil
}

func buildBody(b model.BodySpec) ([]byte, string, error) {
	switch model.Encoding(b.Encoding) {
	case model.EncodingNone, "":
		return nil, "", nil
	case model.EncodingJSON:
		if len(b.JSON) == 0 {
			return []byte("null"), "application/json", nil
		}
		if !json.Valid(b.JSON) {
			return nil, "", errors.New("body.json is not valid JSON")
		}
		return []byte(b.JSON), "application/json", nil
	case model.EncodingForm:
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
	case model.EncodingRaw:
		return []byte(b.Raw), "", nil
	case model.EncodingMultipart:
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
