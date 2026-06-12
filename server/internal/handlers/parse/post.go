package parse

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	parseType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/parse"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
	"github.com/AdityaSinghRajawat/tryit/server/internal/validations"
)

func (h *ParseHandler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req parseType.Request
	if err := utils.DecodeJSONRequest(r, &req); err != nil {
		utils.HandleCustomError(
			w,
			utils.LogAndReturnCustomErr(ctx, err, config.GetErrCodeInvalidRequest()),
		)
		return
	}
	if err := validations.ValidateStruct(req); err != nil {
		utils.HandleCustomError(
			w,
			utils.LogAndReturnCustomErr(ctx, err, config.GetErrCodeInvalidRequest()),
		)
		return
	}

	resp, customErr := h.Service.Parse(ctx, req)
	if customErr != nil {
		utils.HandleCustomError(w, customErr)
		return
	}

	utils.BuildAndSendResponse(ctx, w, resp, http.StatusOK)
}
