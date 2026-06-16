package execute

import (
	"context"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	executeType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/execute"
	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (s *ExecuteService) ExecuteCommand(
	ctx context.Context,
	spec specType.RequestSpec,
	refs map[string]string,
) (*executeType.Response, *config.CustomError) {
	req, err := utils.BuildHTTPRequest(ctx, spec)
	if err != nil {
		return nil, config.NewCustomError(err, config.GetErrCodeInvalidRequest())
	}

	if need := s.checkConsent(req, spec, refs); need != nil {
		return &executeType.Response{
			RequestPreview: executeType.RequestPreview{
				Method:  req.Method,
				URL:     req.URL.String(),
				Headers: utils.FlattenHeaders(req.Header),
				Body:    utils.RequestBodyPreview(req),
			},
			ConsentRequired: need,
		}, nil
	}

	preview := executeType.RequestPreview{
		Method:  req.Method,
		URL:     req.URL.String(),
		Headers: utils.FlattenHeaders(req.Header),
		Body:    utils.RequestBodyPreview(req),
	}

	maskedAuth, cerr := s.injectAuth(req, spec, refs)
	if cerr != nil {
		return nil, cerr
	}
	if maskedAuth != "" {
		preview.Headers[config.GetHeaderAuthorization()] = utils.StripHeaderName(maskedAuth)
	}

	start := utils.GetCurrTimeStamp()
	res, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, config.NewCustomError(err, config.GetErrCodeTargetUnreachable())
	}

	return &executeType.Response{
		Status:          res.Status,
		DurationMs:      utils.ComputeTimeTaken(start),
		ResponseHeaders: utils.FlattenHeaders(res.Headers),
		Body:            string(res.Body),
		Truncated:       res.Truncated,
		RequestPreview:  preview,
	}, nil
}
