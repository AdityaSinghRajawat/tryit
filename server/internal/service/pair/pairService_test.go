package pair

import (
	"errors"
	"testing"
)

type memStore struct {
	token string
	bound string
}

func (m *memStore) Token() string                       { return m.token }
func (m *memStore) BoundOrigin() string                 { return m.bound }
func (m *memStore) SetBoundOrigin(o string) error       { m.bound = o; return nil }
func (m *memStore) Reset() error                        { m.bound = ""; return nil }

func TestVerifyBindsOriginOnFirstMatch(t *testing.T) {
	s := New(&memStore{token: "tk_x"})
	got, err := s.Verify("tk_x", "chrome-extension://abcd")
	if err != nil || got != "chrome-extension://abcd" {
		t.Fatalf("bind: got=%q err=%v", got, err)
	}
	if s.BoundOrigin() != "chrome-extension://abcd" {
		t.Fatalf("BoundOrigin not persisted")
	}
}

func TestVerifyRejectsBadToken(t *testing.T) {
	s := New(&memStore{token: "tk_x"})
	_, err := s.Verify("tk_wrong", "chrome-extension://abcd")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("want ErrInvalidToken, got %v", err)
	}
}

func TestVerifyRejectsNonExtensionOrigin(t *testing.T) {
	s := New(&memStore{token: "tk_x"})
	_, err := s.Verify("tk_x", "https://docs.example.com")
	if !errors.Is(err, ErrInvalidOrigin) {
		t.Fatalf("want ErrInvalidOrigin, got %v", err)
	}
}

func TestVerifyRejectsDifferentOriginOnceBound(t *testing.T) {
	s := New(&memStore{token: "tk_x", bound: "chrome-extension://A"})
	_, err := s.Verify("tk_x", "chrome-extension://B")
	if !errors.Is(err, ErrOriginConflict) {
		t.Fatalf("want ErrOriginConflict, got %v", err)
	}
}

func TestVerifyAcceptsRepeatPair(t *testing.T) {
	s := New(&memStore{token: "tk_x", bound: "chrome-extension://A"})
	got, err := s.Verify("tk_x", "chrome-extension://A")
	if err != nil || got != "chrome-extension://A" {
		t.Fatalf("repeat: got=%q err=%v", got, err)
	}
}
