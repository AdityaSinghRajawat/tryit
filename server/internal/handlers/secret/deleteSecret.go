package secret

import (
	"net/http"

	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
	"github.com/go-chi/chi/v5"
)

func (h *SecretHandler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := chi.URLParam(r, "name")
	if cerr := h.SecretService.DeleteSecret(name); cerr != nil {
		utils.HandleCustomError(w, cerr)
		return
	}

	utils.BuildAndSendResponse(
		ctx,
		w,
		secretType.DeleteResponse{Name: name, Deleted: true},
		http.StatusOK,
	)
}
