package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// errorCodeToStatus maps every error code to its HTTP status (IMPL §9.2).
// Add new entries when new codes appear.
var errorCodeToStatus = map[config.CustomErrorCode]int{
	config.GetErrCodeUnauthorized():      http.StatusUnauthorized,
	config.GetErrCodeForbiddenOrigin():   http.StatusForbidden,
	config.GetErrCodeNotPaired():         http.StatusConflict,
	config.GetErrCodeInvalidRequest():    http.StatusBadRequest,
	config.GetErrCodeParseFailed():       http.StatusUnprocessableEntity,
	config.GetErrCodeSecretNotFound():    http.StatusNotFound,
	config.GetErrCodeAIUnavailable():     http.StatusServiceUnavailable,
	config.GetErrCodeTargetUnreachable(): http.StatusBadGateway,
	config.GetErrCodeInternal():          http.StatusInternalServerError,
}

// errorEnvelope is the wire shape returned to the panel for non-data errors.
type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    config.CustomErrorCode `json:"code"`
	Message string                 `json:"message"`
}

// LogAndReturnCustomErr is the standard way services produce a *CustomError:
// it logs (Phase 2 plugs a context-aware logger here) and returns the typed
// error so handlers can pass it straight to HandleCustomError.
func LogAndReturnCustomErr(
	_ context.Context,
	err error,
	code config.CustomErrorCode,
) *config.CustomError {
	// TODO(phase 2): pull a slog.Logger from ctx and emit an error line.
	return config.NewCustomError(err, code)
}

// HandleCustomError writes a *CustomError as the standard envelope, mapping
// its code to the right HTTP status. Safe with a nil customErr (no-op).
func HandleCustomError(w http.ResponseWriter, customErr *config.CustomError) {
	if customErr == nil {
		return
	}
	status, ok := errorCodeToStatus[customErr.ErrCode]
	if !ok {
		status = http.StatusInternalServerError
	}
	msg := "unknown error"
	if customErr.Error != nil {
		msg = customErr.Error.Error()
	}
	writeJSON(w, status, errorEnvelope{Error: errorBody{Code: customErr.ErrCode, Message: msg}})
}

// BuildAndSendResponse writes a success response. If marshaling fails, it
// emits an internal error envelope.
func BuildAndSendResponse(ctx context.Context, w http.ResponseWriter, resp any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		HandleCustomError(w, LogAndReturnCustomErr(ctx, err, config.GetErrCodeInternal()))
	}
}

// DecodeJSONRequest decodes a JSON body strictly: unknown fields are rejected.
// Use this in every handler that takes a typed request body.
func DecodeJSONRequest(r *http.Request, v any) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	return d.Decode(v)
}

// writeJSON is internal; callers should use BuildAndSendResponse for success
// and HandleCustomError for failures so the wire shape is consistent.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
