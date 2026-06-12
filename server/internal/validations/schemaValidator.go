// Package validations: JSON-Schema + validator/v10 surface. Custom field
// validators land in their own file per domain and self-register in init().
package validations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type SchemaValidator struct {
	schema *jsonschema.Schema
	raw    []byte
}

func NewSchemaValidator(schemaBytes []byte) (*SchemaValidator, error) {
	if len(schemaBytes) == 0 {
		return nil, errors.New("schema bytes are empty")
	}
	var raw any
	if err := json.Unmarshal(schemaBytes, &raw); err != nil {
		return nil, fmt.Errorf("schema is not valid JSON: %w", err)
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource("file:///schema.json", raw); err != nil {
		return nil, fmt.Errorf("add schema resource: %w", err)
	}
	s, err := c.Compile("file:///schema.json")
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
	}
	return &SchemaValidator{schema: s, raw: schemaBytes}, nil
}

// Validate returns an error whose text is suitable for an AI repair-retry prompt.
func (v *SchemaValidator) Validate(payload []byte) error {
	if v == nil || v.schema == nil {
		return errors.New("schema validator not initialised")
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var doc any
	if err := dec.Decode(&doc); err != nil {
		return fmt.Errorf("payload is not valid JSON: %w", err)
	}
	return v.schema.Validate(doc)
}

// Raw returns the original schema bytes for prompt embedding.
func (v *SchemaValidator) Raw() []byte {
	if v == nil {
		return nil
	}
	return v.raw
}
