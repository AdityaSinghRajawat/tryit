// Package pair serves POST /pair (IMPL §8.1). Exempt from the security
// middleware — the body token is the auth; the pair service performs its own
// constant-time compare.
package pair

import pairSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/pair"

type PairHandler struct {
	Service *pairSvc.PairService
}

func NewPairHandler(svc *pairSvc.PairService) *PairHandler {
	return &PairHandler{Service: svc}
}
