package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/dto"
	"github.com/AdityaSinghRajawat/tryit/server/internal/service/execute"
	"github.com/AdityaSinghRajawat/tryit/server/internal/service/secret"
)

type ExecuteHandler struct {
	svc *execute.Service
}

func NewExecuteHandler(svc *execute.Service) *ExecuteHandler { return &ExecuteHandler{svc: svc} }

func (h *ExecuteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, dto.ErrInvalidRequest, "POST required")
		return
	}
	var req dto.ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, dto.ErrInvalidRequest, "invalid JSON body")
		return
	}
	if err := req.RequestSpec.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, dto.ErrInvalidRequest, err.Error())
		return
	}

	out, err := h.svc.Execute(r.Context(), req.RequestSpec, req.SecretRefs)
	if err != nil {
		switch {
		case errors.Is(err, secret.ErrNotFound):
			writeError(w, http.StatusNotFound, dto.ErrSecretNotFound, err.Error())
		default:
			// Distinguish network-y errors from internal: a context-deadline
			// or transport error from the target is "target_unreachable".
			// We treat anything that escaped the service layer as a transport
			// failure for Phase 1; the panel surfaces a useful message.
			writeError(w, http.StatusBadGateway, dto.ErrTargetUnreachable, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, dto.ExecuteResponse{
		Status:          out.Status,
		DurationMs:      out.DurationMs,
		ResponseHeaders: out.Headers,
		Body:            out.Body,
		Truncated:       out.Truncated,
		RequestPreview: dto.RequestPreview{
			Method:  out.RequestPreview.Method,
			URL:     out.RequestPreview.URL,
			Headers: out.RequestPreview.Headers,
			Body:    out.RequestPreview.Body,
		},
	})
}
