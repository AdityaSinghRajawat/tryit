// Package api: embedded wire contract (RequestSpec schema + builtin profiles).
package api

import _ "embed"

//go:embed requestSpec.schema.json
var Schema []byte

//go:embed profiles.json
var BuiltinProfiles []byte
