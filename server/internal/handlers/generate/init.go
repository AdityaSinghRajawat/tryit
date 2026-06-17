// Package generate serves POST /generate — turns a RequestSpec into a
// runnable code snippet (curl / python / javascript).
package generate

import (
	generateSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/generate"
)

type GenerateHandler struct {
	CodegenService *generateSvc.CodegenService
}

func NewGenerateHandler(codegenService *generateSvc.CodegenService) *GenerateHandler {
	return &GenerateHandler{CodegenService: codegenService}
}
