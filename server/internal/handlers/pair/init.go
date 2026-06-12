// Package pair serves POST /pair — exempt from the security middleware
// because the body token is the auth.
package pair

import pairSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/pair"

type PairHandler struct {
	Service *pairSvc.PairService
}

func NewPairHandler(svc *pairSvc.PairService) *PairHandler {
	return &PairHandler{Service: svc}
}
