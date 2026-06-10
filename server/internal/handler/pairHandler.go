package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/dto"
	"github.com/AdityaSinghRajawat/tryit/server/internal/service/pair"
)

// PairHandler implements §8.1. Exempt from security middleware (the body
// token is the auth); CORS allows /pair preflight from any chrome-extension://.
type PairHandler struct {
	svc *pair.Service
}

func NewPairHandler(svc *pair.Service) *PairHandler { return &PairHandler{svc: svc} }

func (h *PairHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, dto.ErrInvalidRequest, "POST required")
		return
	}
	var req dto.PairRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, dto.ErrInvalidRequest, "invalid JSON body")
		return
	}
	origin := r.Header.Get("Origin")
	bound, err := h.svc.Verify(req.Token, origin)
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, dto.PairResponse{OK: true, BoundOrigin: bound})
	case errors.Is(err, pair.ErrInvalidToken):
		writeError(w, http.StatusUnauthorized, dto.ErrUnauthorized, "invalid token")
	case errors.Is(err, pair.ErrInvalidOrigin):
		writeError(w, http.StatusForbidden, dto.ErrForbiddenOrigin, "origin not allowed")
	case errors.Is(err, pair.ErrOriginConflict):
		writeError(w, http.StatusForbidden, dto.ErrForbiddenOrigin, "origin conflicts with bound origin")
	default:
		writeError(w, http.StatusInternalServerError, dto.ErrInternal, err.Error())
	}
}
