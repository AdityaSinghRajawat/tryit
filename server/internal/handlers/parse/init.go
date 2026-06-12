package parse

import parseSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/parse"

type ParseHandler struct {
	Service *parseSvc.ParseService
}

func NewParseHandler(svc *parseSvc.ParseService) *ParseHandler {
	return &ParseHandler{Service: svc}
}
