// Package spec mirrors api/requestSpec.schema.json and extension/src/shared/types.ts.
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
