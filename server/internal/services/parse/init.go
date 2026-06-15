// Package parse runs the detection cascade: cache → structured hint → AI.
package parse

import (
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
	"github.com/AdityaSinghRajawat/tryit/server/internal/validations"
)

type ParseService struct {
	AIProvider      ai.AIProvider
	Cache           *utils.Cache
	SchemaValidator *validations.SchemaValidator
}

func NewParseService(
	aiProvider ai.AIProvider,
	cache *utils.Cache,
	schemaValidator *validations.SchemaValidator,
) *ParseService {
	return &ParseService{
		AIProvider:      aiProvider,
		Cache:           cache,
		SchemaValidator: schemaValidator,
	}
}
