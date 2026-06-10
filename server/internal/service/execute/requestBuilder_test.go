package execute

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/AdityaSinghRajawat/tryit/server/internal/model"
)

func TestBuildRequestPathParams(t *testing.T) {
	spec := model.RequestSpec{
		Method:  "GET",
		BaseURL: "https://api.twilio.com",
		Path:    "/2010/Accounts/{Sid}/Messages.json",
		PathParams: []model.Param{
			{Name: "Sid", Value: "AC123"},
		},
		Auth:       model.AuthSpec{Type: "none"},
		Body:       model.BodySpec{Encoding: "none"},
		Confidence: 1,
	}
	r, err := buildRequest(context.Background(), spec)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	want := "https://api.twilio.com/2010/Accounts/AC123/Messages.json"
	if r.URL.String() != want {
		t.Errorf("URL = %q want %q", r.URL.String(), want)
	}
}

func TestBuildRequestRepeatedQuery(t *testing.T) {
	spec := model.RequestSpec{
		Method:  "GET",
		BaseURL: "https://api.example.com",
		Path:    "/v1/things",
		Query: []model.Param{
			{Name: "expand[]", Values: []string{"a", "b", "c"}},
		},
		Auth:       model.AuthSpec{Type: "none"},
		Body:       model.BodySpec{Encoding: "none"},
		Confidence: 1,
	}
	r, err := buildRequest(context.Background(), spec)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	q := r.URL.Query()
	if got := q["expand[]"]; len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("expand[] = %v", got)
	}
}

func TestBuildRequestFormBody(t *testing.T) {
	spec := model.RequestSpec{
		Method:  "POST",
		BaseURL: "https://api.stripe.com",
		Path:    "/v1/payment_intents",
		Auth:    model.AuthSpec{Type: "none"},
		Body: model.BodySpec{Encoding: "form", Form: []model.Param{
			{Name: "amount", Value: "2000"},
			{Name: "currency", Value: "usd"},
		}},
		Confidence: 1,
	}
	r, err := buildRequest(context.Background(), spec)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
		t.Errorf("Content-Type = %q", ct)
	}
	b, _ := io.ReadAll(r.Body)
	got := string(b)
	if !strings.Contains(got, "amount=2000") || !strings.Contains(got, "currency=usd") {
		t.Errorf("form body = %q", got)
	}
}

func TestBuildRequestJSONBody(t *testing.T) {
	spec := model.RequestSpec{
		Method:     "POST",
		BaseURL:    "https://x.example.com",
		Path:       "/echo",
		Auth:       model.AuthSpec{Type: "none"},
		Body:       model.BodySpec{Encoding: "json", JSON: json.RawMessage(`{"hello":"world"}`)},
		Confidence: 1,
	}
	r, err := buildRequest(context.Background(), spec)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	b, _ := io.ReadAll(r.Body)
	if string(b) != `{"hello":"world"}` {
		t.Errorf("body = %q", b)
	}
}

func TestBuildRequestRelativeBaseURLRejected(t *testing.T) {
	spec := model.RequestSpec{
		Method: "GET", BaseURL: "/x", Path: "/y",
		Auth: model.AuthSpec{Type: "none"}, Body: model.BodySpec{Encoding: "none"},
	}
	if _, err := buildRequest(context.Background(), spec); err == nil {
		t.Fatalf("expected error for relative baseUrl")
	}
}
