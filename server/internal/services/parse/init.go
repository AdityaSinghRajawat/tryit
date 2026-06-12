// Package parse runs the detection cascade: cache → structured hint → AI.
package parse

import (
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
	"github.com/AdityaSinghRajawat/tryit/server/internal/validations"
)

type ParseService struct {
	AI        ai.AIProvider
	Redis     *utils.RedisUtil
	Validator *validations.SchemaValidator
}

func NewParseService(
	aiProvider ai.AIProvider,
	redis *utils.RedisUtil,
	validator *validations.SchemaValidator,
) *ParseService {
	return &ParseService{AI: aiProvider, Redis: redis, Validator: validator}
}
