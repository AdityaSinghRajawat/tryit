// Package parse runs the detection cascade: cache → structured hint → AI,
// with profile overlay on AI output when a host profile matches.
package parse

import (
	"github.com/AdityaSinghRajawat/tryit/server/internal/integrations/ai"
	profileSvc "github.com/AdityaSinghRajawat/tryit/server/internal/services/profile"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
	"github.com/AdityaSinghRajawat/tryit/server/internal/validations"
)

type ParseService struct {
	AIProvider      ai.AIProvider
	Cache           *utils.Cache
	SchemaValidator *validations.SchemaValidator
	ProfileService  *profileSvc.ProfileService
}

func NewParseService(
	aiProvider ai.AIProvider,
	cache *utils.Cache,
	schemaValidator *validations.SchemaValidator,
	profileService *profileSvc.ProfileService,
) *ParseService {
	return &ParseService{
		AIProvider:      aiProvider,
		Cache:           cache,
		SchemaValidator: schemaValidator,
		ProfileService:  profileService,
	}
}
