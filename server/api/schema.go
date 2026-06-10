// Package api holds the wire contract: the RequestSpec JSON Schema (the single
// source of truth) and built-in site profiles, both embedded at build time so
// the binary needs no runtime files.
package api

import _ "embed"

//go:embed requestSpec.schema.json
var Schema []byte

//go:embed profiles.json
var BuiltinProfiles []byte
