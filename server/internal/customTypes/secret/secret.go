// Package secret holds the masking Secret wrapper. Reveal() is the only
// accessor for the underlying credential value and is called solely by
// services/execute/authInjector at the moment of sending an outbound request.
package secret

type Secret struct {
	name  string
	typ   string // "bearer" | "basic" | "apiKey"
	value string
	user  string
	pass  string
}

func New(name, typ, value string) Secret {
	return Secret{name: name, typ: typ, value: value}
}

func NewBasic(name, user, pass string) Secret {
	return Secret{name: name, typ: "basic", user: user, pass: pass}
}

// String redacts intentionally — never log a real secret value.
func (s Secret) String() string { return "Secret(" + s.name + ")" }

func (s Secret) Name() string { return s.name }
func (s Secret) Type() string { return s.typ }

// Reveal returns the underlying credential triple. Use ONLY at the last
// moment before sending the outbound request.
func (s Secret) Reveal() (value, user, pass string) {
	return s.value, s.user, s.pass
}
