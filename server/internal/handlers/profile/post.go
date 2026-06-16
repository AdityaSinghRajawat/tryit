package profile

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	profileType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/profile"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *ProfileHandler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &profileType.SiteProfile{}
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

	if cerr := h.ProfileService.Learn(req); cerr != nil {
		utils.HandleCustomError(w, cerr)
		return
	}

	utils.BuildAndSendResponse(
		ctx,
		w,
		profileType.CreateResponse{Host: req.Host},
		http.StatusOK,
	)
}
