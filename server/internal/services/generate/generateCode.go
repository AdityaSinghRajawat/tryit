package generate

import (
	"bytes"
	"fmt"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	generateType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/generate"
)

func (s *CodegenService) GenerateCode(
	req *generateType.GenerateRequest,
) (*generateType.GenerateResponse, *config.CustomError) {
	model := s.buildRenderModel(req.RequestSpec)

	var buf bytes.Buffer
	name := string(req.Language) + ".tmpl"
	if err := s.templates.ExecuteTemplate(&buf, name, model); err != nil {
		return nil, config.NewCustomError(
			fmt.Errorf("render %s template: %w", req.Language, err),
			config.GetErrCodeInternal(),
		)
	}

	return &generateType.GenerateResponse{
		Language: req.Language,
		Code:     buf.String(),
	}, nil
}
