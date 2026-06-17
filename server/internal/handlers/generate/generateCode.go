package generate

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	generateType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/generate"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *GenerateHandler) GenerateCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &generateType.GenerateRequest{}
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

	resp, cerr := h.CodegenService.GenerateCode(req)
	if cerr != nil {
		utils.HandleCustomError(w, cerr)
		return
	}

	utils.BuildAndSendResponse(ctx, w, resp, http.StatusOK)
}
