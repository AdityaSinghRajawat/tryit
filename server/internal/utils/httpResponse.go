package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
)

// errorCodeToStatus is the single source of truth for code → HTTP status.
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

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    config.CustomErrorCode `json:"code"`
	Message string                 `json:"message"`
}

func LogAndReturnCustomErr(
	_ context.Context,
	err error,
	code config.CustomErrorCode,
) *config.CustomError {
	// TODO: pull a context-aware logger from ctx and emit an error line.
	return config.NewCustomError(err, code)
}

// HandleCustomError renders the standard error envelope. Safe with nil.
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

func BuildAndSendResponse(ctx context.Context, w http.ResponseWriter, resp any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		HandleCustomError(w, LogAndReturnCustomErr(ctx, err, config.GetErrCodeInternal()))
	}
}

// DecodeJSONRequest decodes strictly: unknown fields are rejected.
func DecodeJSONRequest(r *http.Request, v any) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	return d.Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
