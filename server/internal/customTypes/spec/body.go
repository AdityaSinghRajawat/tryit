package spec

import "encoding/json"

type BodySpec struct {
	Encoding    string          `json:"encoding"              validate:"required,oneof=none json form multipart raw"`
	JSON        json.RawMessage `json:"json,omitempty"`
	Form        []Param         `json:"form,omitempty"`
	Raw         string          `json:"raw,omitempty"`
	ContentType string          `json:"contentType,omitempty"`
}
