// Package pair serves POST /pair — exempt from the security middleware
// because the body token is the auth.
package pair

import pairSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/pair"

type PairHandler struct {
	PairService *pairSvc.PairService
}

func NewPairHandler(pairService *pairSvc.PairService) *PairHandler {
	return &PairHandler{PairService: pairService}
}
