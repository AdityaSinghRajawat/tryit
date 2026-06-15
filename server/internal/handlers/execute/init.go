package execute

import executeSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/execute"

type ExecuteHandler struct {
	ExecuteService *executeSvc.ExecuteService
}

func NewExecuteHandler(executeService *executeSvc.ExecuteService) *ExecuteHandler {
	return &ExecuteHandler{ExecuteService: executeService}
}
