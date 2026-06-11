// Package spec is the Go mirror of api/requestSpec.schema.json (the single
// source of truth). Keep in sync with extension/src/shared/types.ts. The
// contract round-trip test (api/contract_test.go) enforces it.
package spec

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
