package spec

type Header struct {
	Name     string `json:"name"               validate:"required,min=1"`
	Value    string `json:"value"`
	Required bool   `json:"required,omitempty"`
}
