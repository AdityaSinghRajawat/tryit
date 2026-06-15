package parse

import parseSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/parse"

type ParseHandler struct {
	ParseService *parseSvc.ParseService
}

func NewParseHandler(parseService *parseSvc.ParseService) *ParseHandler {
	return &ParseHandler{ParseService: parseService}
}
