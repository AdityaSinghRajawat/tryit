package profile

import (
	"net/http"

	profileType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/profile"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	utils.BuildAndSendResponse(
		r.Context(),
		w,
		profileType.ListResponse{Profiles: h.ProfileService.List()},
		http.StatusOK,
	)
}
