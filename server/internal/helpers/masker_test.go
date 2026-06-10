package helpers

import "testing"

func TestMask(t *testing.T) {
	cases := []struct{ in, want string }{
		{"", ""},
		{"abc", "••••"},
		{"abcd", "••••"},
		{"sk_test_12345abcd", "••••abcd"},
	}
	for _, c := range cases {
		if got := Mask(c.in); got != c.want {
			t.Errorf("Mask(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaskBearer(t *testing.T) {
	if got := MaskBearer("Bearer sk_test_12345abcd"); got != "Bearer ••••abcd" {
		t.Errorf("MaskBearer split form: got %q", got)
	}
	if got := MaskBearer("rawvalueabcd"); got != "••••abcd" {
		t.Errorf("MaskBearer no-space form: got %q", got)
	}
}

func TestBasicAuthValue(t *testing.T) {
	if got := BasicAuthValue("user", "pass"); got != "dXNlcjpwYXNz" {
		t.Errorf("BasicAuthValue = %q", got)
	}
}
