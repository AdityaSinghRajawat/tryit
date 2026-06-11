package health

type Response struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Paired  bool   `json:"paired"`
}
