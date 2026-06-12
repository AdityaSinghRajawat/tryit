package execute

import executeSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/execute"

type ExecuteHandler struct {
	Service *executeSvc.ExecuteService
}

func NewExecuteHandler(svc *executeSvc.ExecuteService) *ExecuteHandler {
	return &ExecuteHandler{Service: svc}
}
