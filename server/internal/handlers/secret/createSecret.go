package secret

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	secretType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/secret"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *SecretHandler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &secretType.CreateRequest{}
	if err := utils.DecodeJSONRequest(r, req); err != nil {
		utils.HandleCustomError(
			w,
			utils.LogAndReturnCustomErr(ctx, err, config.GetErrCodeInvalidRequest()),
		)
		return
	}

	if err := req.Validate(); err != nil {
		utils.HandleCustomError(
			w,
			utils.LogAndReturnCustomErr(ctx, err, config.GetErrCodeInvalidRequest()),
		)
		return
	}

	if cerr := h.SecretService.CreateSecret(req.ToRecord()); cerr != nil {
		utils.HandleCustomError(w, cerr)
		return
	}

	utils.BuildAndSendResponse(
		ctx,
		w,
		secretType.CreateResponse{Name: req.Name},
		http.StatusCreated,
	)
}
