package execute

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/AdityaSinghRajawat/tryit/server/internal/model"
)

type fakeResolver struct{ m map[string]model.Secret }

func (f *fakeResolver) Resolve(name string) (model.Secret, error) {
	if s, ok := f.m[name]; ok {
		return s, nil
	}
	return model.Secret{}, errors.New("not found")
}

func newReq(t *testing.T, method, url string) *http.Request {
	t.Helper()
	r, err := http.NewRequestWithContext(context.Background(), method, url, http.NoBody)
	if err != nil {
		t.Fatalf("newReq: %v", err)
	}
	return r
}

func TestInjectBearer(t *testing.T) {
	r := newReq(t, "GET", "https://api.stripe.com/v1/x")
	res := &fakeResolver{m: map[string]model.Secret{
		"STRIPE_KEY": model.NewSecret("STRIPE_KEY", "bearer", "sk_live_12345abcd"),
	}}
	spec := model.RequestSpec{
		Auth: model.AuthSpec{Type: "bearer", ValueRef: "{{secret:STRIPE_KEY}}"},
	}
	masked, err := injectAuth(r, spec, nil, res)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := r.Header.Get("Authorization"); got != "Bearer sk_live_12345abcd" {
		t.Errorf("Authorization = %q", got)
	}
	if !strings.Contains(masked, "••••abcd") {
		t.Errorf("masked preview should contain last4 mask: %q", masked)
	}
	if strings.Contains(masked, "sk_live_12345") {
		t.Errorf("masked preview leaks secret: %q", masked)
	}
}

func TestInjectAPIKeyInQuery(t *testing.T) {
	r := newReq(t, "GET", "https://maps.googleapis.com/maps/api/geocode/json?address=foo")
	res := &fakeResolver{m: map[string]model.Secret{
		"GMAPS": model.NewSecret("GMAPS", "bearer", "AIzaSyD_secret_value"),
	}}
	spec := model.RequestSpec{
		Auth: model.AuthSpec{Type: "apiKey", In: "query", Name: "key", ValueRef: "{{secret:GMAPS}}"},
	}
	masked, err := injectAuth(r, spec, nil, res)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	got := r.URL.Query().Get("key")
	if got != "AIzaSyD_secret_value" {
		t.Errorf("query key = %q", got)
	}
	if !strings.Contains(masked, "••••") || strings.Contains(masked, "AIzaSyD") {
		t.Errorf("masked preview leaks or missing mask: %q", masked)
	}
}

func TestInjectAPIKeyInHeader(t *testing.T) {
	r := newReq(t, "GET", "https://api.example.com/x")
	res := &fakeResolver{m: map[string]model.Secret{
		"X_KEY": model.NewSecret("X_KEY", "bearer", "raw_api_key_value"),
	}}
	spec := model.RequestSpec{
		Auth: model.AuthSpec{Type: "apiKey", In: "header", Name: "X-API-Key", ValueRef: "{{secret:X_KEY}}"},
	}
	if _, err := injectAuth(r, spec, nil, res); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := r.Header.Get("X-API-Key"); got != "raw_api_key_value" {
		t.Errorf("X-API-Key = %q", got)
	}
}

func TestInjectBasicTwoSecrets(t *testing.T) {
	// Twilio-style: two distinct stored secrets, one for SID one for token.
	r := newReq(t, "POST", "https://api.twilio.com/x")
	res := &fakeResolver{m: map[string]model.Secret{
		"TWILIO_SID":        model.NewSecret("TWILIO_SID", "bearer", "ACabcdef"),
		"TWILIO_AUTH_TOKEN": model.NewSecret("TWILIO_AUTH_TOKEN", "bearer", "secret-token-xyz"),
	}}
	spec := model.RequestSpec{
		Auth: model.AuthSpec{
			Type:     "basic",
			Username: "{{secret:TWILIO_SID}}",
			Password: "{{secret:TWILIO_AUTH_TOKEN}}",
		},
	}
	if _, err := injectAuth(r, spec, nil, res); err != nil {
		t.Fatalf("err: %v", err)
	}
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		t.Fatalf("Authorization = %q (want Basic ... )", auth)
	}
	// Decoded form is "ACabcdef:secret-token-xyz" -> base64 QUNhYmNkZWY6c2VjcmV0LXRva2VuLXh5eg==
	if auth != "Basic QUNhYmNkZWY6c2VjcmV0LXRva2VuLXh5eg==" {
		t.Errorf("Basic value = %q", auth)
	}
}

func TestInjectBasicSingleSecret(t *testing.T) {
	// Single stored secret holding both halves (the model.NewBasicSecret form).
	r := newReq(t, "POST", "https://x/y")
	res := &fakeResolver{m: map[string]model.Secret{
		"creds": model.NewBasicSecret("creds", "user", "pw"),
	}}
	spec := model.RequestSpec{
		Auth: model.AuthSpec{
			Type:     "basic",
			Username: "{{secret:CREDS}}",
			Password: "{{secret:CREDS}}",
		},
	}
	// secretRefs maps placeholder NAME -> stored name (the panel does this).
	refs := map[string]string{"CREDS": "creds"}
	if _, err := injectAuth(r, spec, refs, res); err != nil {
		t.Fatalf("err: %v", err)
	}
	if auth := r.Header.Get("Authorization"); auth != "Basic dXNlcjpwdw==" {
		t.Errorf("Authorization = %q", auth)
	}
}

func TestInjectNoneIsNoop(t *testing.T) {
	r := newReq(t, "GET", "https://x/y")
	spec := model.RequestSpec{Auth: model.AuthSpec{Type: "none"}}
	if _, err := injectAuth(r, spec, nil, &fakeResolver{}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if r.Header.Get("Authorization") != "" {
		t.Errorf("unexpected Authorization header")
	}
}
