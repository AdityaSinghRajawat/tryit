package health

type PairReader interface {
	BoundOrigin() string
}

type HealthHandler struct {
	PairReader PairReader
}

func NewHealthHandler(pairReader PairReader) *HealthHandler {
	return &HealthHandler{PairReader: pairReader}
}
