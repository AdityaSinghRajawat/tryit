package secret

import (
	"net/http"

	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *SecretHandler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	infos, cerr := h.SecretService.ListSecrets()
	if cerr != nil {
		utils.HandleCustomError(w, cerr)
		return
	}
	utils.BuildAndSendResponse(
		ctx,
		w,
		secretType.ListResponse{Secrets: infos},
		http.StatusOK,
	)
}
