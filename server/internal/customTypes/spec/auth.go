package spec

type AuthSpec struct {
	Type     string `json:"type"               validate:"required,oneof=none bearer basic apiKey"`
	In       string `json:"in,omitempty"       validate:"omitempty,oneof=header query"`
	Name     string `json:"name,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
	ValueRef string `json:"valueRef,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// BasicHalf selects user vs password when resolving a single basic-typed Secret.
type BasicHalf int

const (
	BasicHalfUser BasicHalf = iota
	BasicHalfPass
)
