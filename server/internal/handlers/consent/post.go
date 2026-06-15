package consent

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	consentType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/consent"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *ConsentHandler) Post(w http.ResponseWriter, r *http.Request) {
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

	if cerr := h.ConsentService.Grant(req.Secret, req.Host); cerr != nil {
		utils.HandleCustomError(w, cerr)
		return
	}
	utils.BuildAndSendResponse(ctx, w, consentType.Response{Granted: true}, http.StatusOK)
}
