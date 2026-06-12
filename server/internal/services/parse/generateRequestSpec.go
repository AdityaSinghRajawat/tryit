package parse

import (
	"context"
	"encoding/json"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	aiType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/ai"
	parseType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/parse"
	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// generateRequestSpec runs the AI step of the cascade. Validator failures on
// previous attempts are appended to the next prompt so the model can repair.
func (s *ParseService) generateRequestSpec(
	ctx context.Context,
	req parseType.Request,
) (*specType.RequestSpec, *config.CustomError) {
	baseUserMsg := s.buildUserMessage(req)
	aiReq := aiType.Request{
		System:      s.buildSystemPrompt(),
		User:        baseUserMsg,
		Schema:      s.Validator.Raw(),
		Temperature: config.GetAITemperature(),
		MaxTokens:   config.GetAIMaxTokens(),
	}

	var lastErr error
	spec, err := utils.ExecuteWithRetry(
		ctx,
		"parse.generateRequestSpec",
		config.GetAIRepairRetries()+1,
		0,
		func(ctx context.Context) (*specType.RequestSpec, error) {
			if lastErr != nil {
				aiReq.User = baseUserMsg +
					"\n\nYOUR PREVIOUS RESPONSE FAILED VALIDATION:\n" +
					lastErr.Error() +
					"\n\nFix the JSON and respond again with ONLY the corrected JSON object."
			}

			resp, err := s.AI.Complete(ctx, aiReq)
			if err != nil {
				return nil, err
			}

			payload := []byte(utils.StripJSONFences(resp.Content))
			if vErr := s.Validator.Validate(payload); vErr != nil {
				lastErr = vErr
				return nil, vErr
			}

			var parsed specType.RequestSpec
			if jerr := json.Unmarshal(payload, &parsed); jerr != nil {
				lastErr = jerr
				return nil, jerr
			}

			return &parsed, nil
		},
	)
	if err != nil {
		return nil, utils.LogAndReturnCustomErr(ctx, err, config.GetErrCodeParseFailed())
	}

	return spec, nil
}
