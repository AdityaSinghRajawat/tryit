package health

import (
	"net/http"

	healthType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/health"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// Version is set at build time via -ldflags.
var Version = "dev"

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	utils.BuildAndSendResponse(r.Context(), w, healthType.Response{
		Status:  "ok",
		Version: Version,
		Paired:  h.PairReader.BoundOrigin() != "",
	}, http.StatusOK)
}
