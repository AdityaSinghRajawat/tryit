// Package model holds the domain types. requestSpec.go is the Go mirror of
// api/requestSpec.schema.json (the single source of truth). Keep all three
// (schema, this file, extension/src/shared/types.ts) in sync — the contract
// round-trip test (api/contract_test.go) enforces it.
package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Method string

const (
	MethodGET     Method = "GET"
	MethodPOST    Method = "POST"
	MethodPUT     Method = "PUT"
	MethodPATCH   Method = "PATCH"
	MethodDELETE  Method = "DELETE"
	MethodHEAD    Method = "HEAD"
	MethodOPTIONS Method = "OPTIONS"
)

type AuthType string

const (
	AuthNone   AuthType = "none"
	AuthBearer AuthType = "bearer"
	AuthBasic  AuthType = "basic"
	AuthAPIKey AuthType = "apiKey"
)

type Encoding string

const (
	EncodingNone      Encoding = "none"
	EncodingJSON      Encoding = "json"
	EncodingForm      Encoding = "form"
	EncodingMultipart Encoding = "multipart"
	EncodingRaw       Encoding = "raw"
)

type RequestSpec struct {
	Method     string   `json:"method"`
	BaseURL    string   `json:"baseUrl"`
	Path       string   `json:"path"`
	PathParams []Param  `json:"pathParams,omitempty"`
	Query      []Param  `json:"query,omitempty"`
	Headers    []Header `json:"headers,omitempty"`
	Auth       AuthSpec `json:"auth"`
	Body       BodySpec `json:"body"`
	Confidence float64  `json:"confidence"`
	Notes      string   `json:"notes,omitempty"`
}

type Param struct {
	Name        string   `json:"name"`
	Value       string   `json:"value,omitempty"`
	Values      []string `json:"values,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Header struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Required bool   `json:"required,omitempty"`
}

type AuthSpec struct {
	Type     string `json:"type"`
	In       string `json:"in,omitempty"`
	Name     string `json:"name,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
	ValueRef string `json:"valueRef,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type BodySpec struct {
	Encoding    string          `json:"encoding"`
	JSON        json.RawMessage `json:"json,omitempty"`
	Form        []Param         `json:"form,omitempty"`
	Raw         string          `json:"raw,omitempty"`
	ContentType string          `json:"contentType,omitempty"`
}

var secretRe = regexp.MustCompile(`\{\{secret:([A-Z0-9_]+)\}\}`)

// SecretRefs scans Auth.* (the only fields where secrets are allowed by the
// contract) and returns the placeholder NAMEs found, deduped, in stable order.
func (s RequestSpec) SecretRefs() []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, v := range []string{s.Auth.ValueRef, s.Auth.Username, s.Auth.Password} {
		for _, m := range secretRe.FindAllStringSubmatch(v, -1) {
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
// JSON Schema validator (utils/validator) is authoritative; this is a fast
// pre-check for handler input.
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
