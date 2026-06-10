package handler

import "net/http"

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Paired  bool   `json:"paired"`
}

type pairingReader interface {
	BoundOrigin() string
}

// Version is set at build time via -ldflags (Phase 3); default "dev".
var Version = "dev"

type HealthHandler struct {
	pair pairingReader
}

func NewHealthHandler(p pairingReader) *HealthHandler { return &HealthHandler{pair: p} }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status:  "ok",
		Version: Version,
		Paired:  h.pair.BoundOrigin() != "",
	})
}
