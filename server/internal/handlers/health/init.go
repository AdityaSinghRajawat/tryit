package health

type PairReader interface {
	BoundOrigin() string
}

type HealthHandler struct {
	Pair PairReader
}

func NewHealthHandler(pair PairReader) *HealthHandler {
	return &HealthHandler{Pair: pair}
}
