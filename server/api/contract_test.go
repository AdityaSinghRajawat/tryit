package api_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/AdityaSinghRajawat/tryit/server/api"
	"github.com/AdityaSinghRajawat/tryit/server/internal/model"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// TestContractRoundTrip enforces D4: the JSON Schema, the Go struct, and (by
// hand-maintenance) the TS type stay in sync. We exercise:
//   - schema parses
//   - a representative RequestSpec marshals to JSON
//   - that JSON validates against the schema
//   - that JSON unmarshals back into the Go struct without information loss
func TestContractRoundTrip(t *testing.T) {
	c := jsonschema.NewCompiler()
	if err := c.AddResource("file:///requestSpec.schema.json", mustJSON(t, api.Schema)); err != nil {
		t.Fatalf("add resource: %v", err)
	}
	sch, err := c.Compile("file:///requestSpec.schema.json")
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}

	cases := map[string]model.RequestSpec{
		"bearer-form": {
			Method: "POST", BaseURL: "https://api.stripe.com", Path: "/v1/payment_intents",
			Auth: model.AuthSpec{Type: "bearer", In: "header", Name: "Authorization", Prefix: "Bearer ",
				ValueRef: "{{secret:STRIPE_API_KEY}}"},
			Body: model.BodySpec{Encoding: "form", Form: []model.Param{
				{Name: "amount", Value: "2000", Required: true},
				{Name: "currency", Value: "usd", Required: true},
			}},
			Confidence: 0.94,
		},
		"basic-form-path": {
			Method:  "POST",
			BaseURL: "https://api.twilio.com",
			Path:    "/2010-04-01/Accounts/{AccountSid}/Messages.json",
			PathParams: []model.Param{
				{Name: "AccountSid", Value: "{{secret:TWILIO_SID}}", Required: true},
			},
			Auth: model.AuthSpec{Type: "basic", Username: "{{secret:TWILIO_SID}}", Password: "{{secret:TWILIO_AUTH_TOKEN}}"},
			Body: model.BodySpec{Encoding: "form", Form: []model.Param{
				{Name: "To", Required: true}, {Name: "From", Required: true}, {Name: "Body", Required: true},
			}},
			Confidence: 0.9,
		},
		"apikey-query": {
			Method:  "GET",
			BaseURL: "https://maps.googleapis.com",
			Path:    "/maps/api/geocode/json",
			Query: []model.Param{
				{Name: "address", Value: "1600 Amphitheatre Parkway", Required: true},
			},
			Auth:       model.AuthSpec{Type: "apiKey", In: "query", Name: "key", ValueRef: "{{secret:GOOGLE_MAPS_KEY}}"},
			Body:       model.BodySpec{Encoding: "none"},
			Confidence: 0.92,
		},
		"none-auth-json": {
			Method:     "POST",
			BaseURL:    "https://example.com",
			Path:       "/echo",
			Auth:       model.AuthSpec{Type: "none"},
			Body:       model.BodySpec{Encoding: "json", JSON: json.RawMessage(`{"hello":"world"}`)},
			Confidence: 1,
		},
	}

	for name, spec := range cases {
		t.Run(name, func(t *testing.T) {
			if err := spec.Validate(); err != nil {
				t.Fatalf("model.Validate: %v", err)
			}
			b, err := json.Marshal(spec)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var v any
			if err := json.Unmarshal(b, &v); err != nil {
				t.Fatalf("unmarshal to any: %v", err)
			}
			if err := sch.Validate(v); err != nil {
				t.Fatalf("schema validate: %v\nspec=%s", err, b)
			}
			var back model.RequestSpec
			if err := json.Unmarshal(b, &back); err != nil {
				t.Fatalf("unmarshal back: %v", err)
			}
			b2, _ := json.Marshal(back)
			if !bytes.Equal(b, b2) {
				t.Fatalf("round-trip drift\n  before: %s\n  after:  %s", b, b2)
			}
		})
	}
}

func TestSecretRefs(t *testing.T) {
	s := model.RequestSpec{
		Auth: model.AuthSpec{
			Type:     "bearer",
			ValueRef: "Bearer {{secret:STRIPE_KEY}}",
			Username: "{{secret:TWILIO_SID}}",
			Password: "{{secret:TWILIO_AUTH_TOKEN}}",
		},
	}
	got := s.SecretRefs()
	want := []string{"STRIPE_KEY", "TWILIO_SID", "TWILIO_AUTH_TOKEN"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestValidateRejectsBadAuth(t *testing.T) {
	bad := model.RequestSpec{
		Method: "GET", BaseURL: "https://x", Path: "/", Confidence: 0.5,
		Auth: model.AuthSpec{Type: "bearer"}, // missing valueRef
		Body: model.BodySpec{Encoding: "none"},
	}
	if err := bad.Validate(); err == nil {
		t.Fatalf("expected error for bearer without valueRef")
	}
}

func mustJSON(t *testing.T, b []byte) any {
	t.Helper()
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("schema is not valid JSON: %v", err)
	}
	return v
}
