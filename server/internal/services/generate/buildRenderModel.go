package generate

import (
	"strings"

	generateType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/generate"
	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

func (s *CodegenService) buildRenderModel(spec specType.RequestSpec) generateType.RenderModel {
	m := generateType.RenderModel{
		Method:      strings.ToUpper(spec.Method),
		BodyEnc:     spec.Body.Encoding,
		ContentType: spec.Body.ContentType,
	}

	pathWithParams := utils.SubstitutePathParams(spec.Path, spec.PathParams)
	queryStr := utils.BuildQueryString(spec.Query)
	m.QueryString = queryStr

	m.URL = strings.TrimRight(spec.BaseURL, "/") + pathWithParams
	if queryStr != "" {
		m.URL += "?" + queryStr
	}

	for _, h := range spec.Headers {
		m.Headers = append(m.Headers, generateType.KV{Name: h.Name, Value: h.Value})
	}

	switch specType.Encoding(spec.Body.Encoding) {
	case specType.EncodingJSON:
		if len(spec.Body.JSON) > 0 {
			m.BodyJSON = utils.PrettyJSON(spec.Body.JSON)
		}
		if m.ContentType == "" {
			m.ContentType = "application/json"
		}
	case specType.EncodingForm:
		for _, p := range spec.Body.Form {
			m.BodyForm = append(m.BodyForm, generateType.KV{Name: p.Name, Value: p.Value})
		}
		if m.ContentType == "" {
			m.ContentType = "application/x-www-form-urlencoded"
		}
	case specType.EncodingMultipart:
		for _, p := range spec.Body.Form {
			m.BodyForm = append(m.BodyForm, generateType.KV{Name: p.Name, Value: p.Value})
		}
		if m.ContentType == "" {
			m.ContentType = "multipart/form-data"
		}
	case specType.EncodingRaw:
		m.BodyRaw = spec.Body.Raw
	}

	s.applyAuth(&m, spec.Auth)
	return m
}

func (s *CodegenService) applyAuth(m *generateType.RenderModel, auth specType.AuthSpec) {
	m.AuthType = auth.Type
	switch specType.AuthType(auth.Type) {
	case specType.AuthBearer:
		envName, _ := utils.PlaceholderToEnv(auth.ValueRef)
		if envName == "" {
			return
		}
		headerName := auth.Name
		if headerName == "" {
			headerName = "Authorization"
		}
		prefix := auth.Prefix
		if prefix == "" {
			prefix = "Bearer "
		}
		m.AuthHeaderName = headerName
		m.AuthHeaderPrefix = prefix
		m.AuthHeaderEnv = envName
		m.AuthHeaderValue = prefix + "$" + envName
		m.EnvVars = utils.AppendUnique(m.EnvVars, envName)

	case specType.AuthBasic:
		userEnv, _ := utils.PlaceholderToEnv(auth.Username)
		passEnv, _ := utils.PlaceholderToEnv(auth.Password)
		m.BasicUserEnv = userEnv
		m.BasicPassEnv = passEnv
		if userEnv != "" {
			m.EnvVars = utils.AppendUnique(m.EnvVars, userEnv)
		}
		if passEnv != "" {
			m.EnvVars = utils.AppendUnique(m.EnvVars, passEnv)
		}

	case specType.AuthAPIKey:
		envName, _ := utils.PlaceholderToEnv(auth.ValueRef)
		if envName == "" {
			return
		}
		m.APIKeyIn = auth.In
		m.APIKeyName = auth.Name
		m.APIKeyEnv = envName
		m.EnvVars = utils.AppendUnique(m.EnvVars, envName)
	}
}
