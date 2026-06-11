package spec

import (
	"fmt"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

type RequestSpec struct {
	Method     string   `json:"method"               validate:"required,oneof=GET POST PUT PATCH DELETE HEAD OPTIONS"`
	BaseURL    string   `json:"baseUrl"              validate:"required,url"`
	Path       string   `json:"path"                 validate:"required"`
	PathParams []Param  `json:"pathParams,omitempty" validate:"dive"`
	Query      []Param  `json:"query,omitempty"      validate:"dive"`
	Headers    []Header `json:"headers,omitempty"    validate:"dive"`
	Auth       AuthSpec `json:"auth"                 validate:"required"`
	Body       BodySpec `json:"body"                 validate:"required"`
	Confidence float64  `json:"confidence"           validate:"gte=0,lte=1"`
	Notes      string   `json:"notes,omitempty"`
}

// SecretRefs scans Auth.* (the only fields where secrets are allowed by the
// contract) and returns the placeholder NAMEs found, deduped, in stable order.
func (s RequestSpec) SecretRefs() []string {
	re := config.GetSecretPlaceholderRegex()
	seen := map[string]struct{}{}
	out := []string{}
	for _, v := range []string{s.Auth.ValueRef, s.Auth.Username, s.Auth.Password} {
		for _, m := range re.FindAllStringSubmatch(v, -1) {
			if _, ok := seen[m[1]]; ok {
				continue
			}
			seen[m[1]] = struct{}{}
			out = append(out, m[1])
		}
	}

	return out
}

// Validate mirrors the schema's allOf rules and enum constraints. The full
// JSON Schema validator (utils) is authoritative; this is a fast pre-check
// for handler input.
func (s RequestSpec) Validate() error {
	switch Method(strings.ToUpper(s.Method)) {
	case MethodGET, MethodPOST, MethodPUT, MethodPATCH, MethodDELETE, MethodHEAD, MethodOPTIONS:
	default:
		return fmt.Errorf("invalid method %q", s.Method)
	}
	if s.BaseURL == "" {
		return fmt.Errorf("baseUrl is required")
	}
	if strings.HasSuffix(s.BaseURL, "/") {
		return fmt.Errorf("baseUrl must not end with a slash")
	}
	if s.Path == "" {
		return fmt.Errorf("path is required")
	}
	if s.Confidence < 0 || s.Confidence > 1 {
		return fmt.Errorf("confidence must be in [0,1], got %v", s.Confidence)
	}
	switch AuthType(s.Auth.Type) {
	case AuthNone:
	case AuthBearer:
		if s.Auth.ValueRef == "" {
			return fmt.Errorf("bearer auth requires valueRef")
		}
	case AuthAPIKey:
		if s.Auth.In == "" || s.Auth.Name == "" || s.Auth.ValueRef == "" {
			return fmt.Errorf("apiKey auth requires in, name, valueRef")
		}
	case AuthBasic:
		if s.Auth.Username == "" || s.Auth.Password == "" {
			return fmt.Errorf("basic auth requires username and password")
		}
	default:
		return fmt.Errorf("invalid auth.type %q", s.Auth.Type)
	}
	switch Encoding(s.Body.Encoding) {
	case EncodingNone, EncodingJSON, EncodingForm, EncodingMultipart, EncodingRaw:
	default:
		return fmt.Errorf("invalid body.encoding %q", s.Body.Encoding)
	}

	return nil
}
