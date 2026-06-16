package consent

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	consentType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/consent"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *ConsentHandler) CreateConsent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &consentType.Request{}
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

	if customErr := h.ConsentService.GrantConsent(req.Secret, req.Host); customErr != nil {
		utils.HandleCustomError(w, customErr)
		return
	}
	utils.BuildAndSendResponse(ctx, w, consentType.Response{Granted: true}, http.StatusOK)
}
