package parse

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	parseType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/parse"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *ParseHandler) ParseCommand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &parseType.Request{}
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

	resp, customErr := h.ParseService.ParseCommand(ctx, *req)
	if customErr != nil {
		utils.HandleCustomError(w, customErr)
		return
	}

	utils.BuildAndSendResponse(ctx, w, resp, http.StatusOK)
}
