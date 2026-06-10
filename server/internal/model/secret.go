package model

// Secret wraps a sensitive value behind Reveal(). Its String/MarshalJSON/etc.
// must never expose the plaintext — Reveal is the only accessor, used solely
// by service/execute/authInjector at the moment of building an outbound
// request. Storage layers MUST construct Secret values via NewSecret; nothing
// else may.
type Secret struct {
	name string
	typ  string // "bearer" | "basic" | "apiKey"
	// For bearer/apiKey: value holds the secret; user/pass empty.
	// For basic: user/pass hold the credentials; value empty.
	value string
	user  string
	pass  string
}

func NewSecret(name, typ, value string) Secret {
	return Secret{name: name, typ: typ, value: value}
}

func NewBasicSecret(name, user, pass string) Secret {
	return Secret{name: name, typ: "basic", user: user, pass: pass}
}

// String redacts intentionally — never log a real secret value.
func (s Secret) String() string { return "Secret(" + s.name + ")" }

func (s Secret) Name() string { return s.name }
func (s Secret) Type() string { return s.typ }

// Reveal returns the underlying credential triple. Use ONLY at the last
// moment before sending the outbound request (service/execute/authInjector).
func (s Secret) Reveal() (value, user, pass string) {
	return s.value, s.user, s.pass
}
