package secret

// CreateRequest covers all auth types — bearer/apiKey use Value; basic uses
// Username + Password.
type CreateRequest struct {
	Name     string `json:"name"               validate:"required,min=1"`
	Type     string `json:"type"               validate:"required,oneof=bearer basic apiKey"`
	HostHint string `json:"hostHint,omitempty"`
	Value    string `json:"value,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type CreateResponse struct {
	Name string `json:"name"`
}

type DeleteResponse struct {
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
}
