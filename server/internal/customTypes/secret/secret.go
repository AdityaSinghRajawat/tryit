// Package secret holds the masking Secret wrapper. Reveal() is the only
// underlying-value accessor; call it ONLY at the moment of sending.
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

// String intentionally redacts — never log a Secret value directly.
func (s Secret) String() string { return "Secret(" + s.name + ")" }

func (s Secret) Name() string { return s.name }
func (s Secret) Type() string { return s.typ }

func (s Secret) Reveal() (value, user, pass string) {
	return s.value, s.user, s.pass
}
