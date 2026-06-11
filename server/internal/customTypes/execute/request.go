package execute

import specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"

type Request struct {
	RequestSpec specType.RequestSpec `json:"requestSpec"          validate:"required"`
	SecretRefs  map[string]string    `json:"secretRefs,omitempty"`
}
