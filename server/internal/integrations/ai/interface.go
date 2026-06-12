package ai

import (
	"context"

	aiType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/ai"
)

// AIProvider is the Strategy contract every provider implementation satisfies.
type AIProvider interface {
	Complete(ctx context.Context, req aiType.Request) (*aiType.Response, error)
}
