package profile

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	profileType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/profile"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
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

	if customErr := h.ProfileService.LearnProfile(req); customErr != nil {
		utils.HandleCustomError(w, customErr)
		return
	}

	utils.BuildAndSendResponse(
		ctx,
		w,
		profileType.CreateResponse{Host: req.Host},
		http.StatusOK,
	)
}
