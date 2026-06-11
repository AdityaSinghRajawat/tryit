package execute

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	executeType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/execute"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *ExecuteHandler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req executeType.Request
	if err := utils.DecodeJSONRequest(r, &req); err != nil {
		utils.HandleCustomError(
			w,
			utils.LogAndReturnCustomErr(ctx, err, config.GetErrCodeInvalidRequest()),
		)
		return
	}
	if err := req.RequestSpec.Validate(); err != nil {
		utils.HandleCustomError(
			w,
			utils.LogAndReturnCustomErr(ctx, err, config.GetErrCodeInvalidRequest()),
		)
		return
	}

	resp, customErr := h.Service.Execute(ctx, req.RequestSpec, req.SecretRefs)
	if customErr != nil {
		utils.HandleCustomError(w, customErr)
		return
	}

	utils.BuildAndSendResponse(ctx, w, resp, http.StatusOK)
}
