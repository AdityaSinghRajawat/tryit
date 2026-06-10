package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/dto"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code dto.ErrorCode, msg string) {
	writeJSON(w, status, dto.ErrorEnvelope{Error: dto.ErrorBody{Code: code, Message: msg}})
}
