package pair

type Request struct {
	Token string `json:"token" validate:"required,min=1"`
}
