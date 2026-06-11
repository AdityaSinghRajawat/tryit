package spec

type Param struct {
	Name        string   `json:"name"                  validate:"required,min=1"`
	Value       string   `json:"value,omitempty"`
	Values      []string `json:"values,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Description string   `json:"description,omitempty"`
}
