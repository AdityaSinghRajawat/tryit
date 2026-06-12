package consent

type Request struct {
	Secret string `json:"secret" validate:"required,min=1"`
	Host   string `json:"host"   validate:"required,min=1"`
}
