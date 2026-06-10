package dto

import "github.com/AdityaSinghRajawat/tryit/server/internal/model"

type ExecuteRequest struct {
	RequestSpec model.RequestSpec `json:"requestSpec"`
	SecretRefs  map[string]string `json:"secretRefs,omitempty"`
}

type RequestPreview struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body,omitempty"`
}

type ExecuteResponse struct {
	Status          int               `json:"status"`
	DurationMs      int64             `json:"durationMs"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	Body            string            `json:"body"`
	Truncated       bool              `json:"truncated"`
	RequestPreview  RequestPreview    `json:"requestPreview"`
	ConsentRequired *ConsentRequired  `json:"consentRequired,omitempty"`
}

type ConsentRequired struct {
	Secret string `json:"secret"`
	Host   string `json:"host"`
}
