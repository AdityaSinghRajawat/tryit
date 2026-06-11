package execute

type ConsentRequired struct {
	Secret string `json:"secret"`
	Host   string `json:"host"`
}

type Response struct {
	Status          int               `json:"status"`
	DurationMs      int64             `json:"durationMs"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	Body            string            `json:"body"`
	Truncated       bool              `json:"truncated"`
	RequestPreview  RequestPreview    `json:"requestPreview"`
	ConsentRequired *ConsentRequired  `json:"consentRequired,omitempty"`
}
