package pair

import (
	"net/http"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	pairType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/pair"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *PairHandler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &pairType.Request{}
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

	origin := r.Header.Get(config.GetHeaderOrigin())
	bound, customErr := h.PairService.Verify(req.Token, origin)
	if customErr != nil {
		utils.HandleCustomError(w, customErr)
		return
	}

	utils.BuildAndSendResponse(
		ctx,
		w,
		pairType.Response{OK: true, BoundOrigin: bound},
		http.StatusOK,
	)
}
