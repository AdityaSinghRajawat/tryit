// Package health serves GET /health (route-exempt from auth).
package health

// PairReader is the narrow surface the health handler needs from the pair
// store to compute the "paired" flag in the response.
type PairReader interface {
	BoundOrigin() string
}

type HealthHandler struct {
	Pair PairReader
}

func NewHealthHandler(pair PairReader) *HealthHandler {
	return &HealthHandler{Pair: pair}
}
