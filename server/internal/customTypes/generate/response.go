package generate

type GenerateResponse struct {
	Language Language `json:"language"`
	Code     string   `json:"code"`
}
