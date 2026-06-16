package execute

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/config"
	specType "github.com/AdityaSinghRajawat/tryit/server/internal/customTypes/spec"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils"
)

// injectAuth resolves {{secret:NAME}} placeholders and stamps the resulting
// header / query on req. Returns the masked preview line for the request
// preview — never the real secret.
func (s *ExecuteService) injectAuth(
	req *http.Request,
	spec specType.RequestSpec,
	mappedRefs map[string]string,
) (string, *config.CustomError) {
	a := spec.Auth

	switch specType.AuthType(a.Type) {
	case specType.AuthNone, "":
		return "", nil

	case specType.AuthBearer:
		val, cerr := s.resolveScalar(a.ValueRef, mappedRefs)
		if cerr != nil {
			return "", cerr
		}
		header := headerOrDefault(a.Name, config.GetHeaderAuthorization())
		prefix := a.Prefix
		if prefix == "" {
			prefix = config.GetPrefixBearer()
		}
		req.Header.Set(header, prefix+val)
		return header + ": " + utils.MaskBearer(prefix+val), nil

	case specType.AuthAPIKey:
		val, cerr := s.resolveScalar(a.ValueRef, mappedRefs)
		if cerr != nil {
			return "", cerr
		}
		header := headerOrDefault(a.Name, "")
		if header == "" {
			return "", config.NewCustomError(
				errors.New("apiKey auth requires a name"),
				config.GetErrCodeInvalidRequest(),
			)
		}
		switch a.In {
		case "query":
			u := req.URL
			q := u.Query()
			q.Set(header, a.Prefix+val)
			u.RawQuery = q.Encode()
			return "?" + header + "=" + utils.Mask(val), nil
		case "header", "":
			req.Header.Set(header, a.Prefix+val)
			return header + ": " + utils.Mask(val), nil
		default:
			return "", config.NewCustomError(
				fmt.Errorf("invalid apiKey auth.in %q", a.In),
				config.GetErrCodeInvalidRequest(),
			)
		}

	case specType.AuthBasic:
		user, cerr := s.resolveBasicHalf(a.Username, mappedRefs, specType.BasicHalfUser)
		if cerr != nil {
			return "", cerr
		}
		pass, cerr := s.resolveBasicHalf(a.Password, mappedRefs, specType.BasicHalfPass)
		if cerr != nil {
			return "", cerr
		}
		req.Header.Set(
			config.GetHeaderAuthorization(),
			config.GetPrefixBasic()+utils.BasicAuthValue(user, pass),
		)
		return config.GetHeaderAuthorization() + ": " + config.GetPrefixBasic() + "••••", nil

	default:
		return "", config.NewCustomError(
			fmt.Errorf("unknown auth.type %q", a.Type),
			config.GetErrCodeInvalidRequest(),
		)
	}
}

func (s *ExecuteService) resolveScalar(
	template string,
	mapped map[string]string,
) (string, *config.CustomError) {
	name, cerr := parseSecretRef(template, mapped)
	if cerr != nil {
		return "", cerr
	}
	sec, cerr := s.SecretService.ResolveSecret(name)
	if cerr != nil {
		return "", cerr
	}
	v, _, _ := sec.Reveal()
	return v, nil
}

func (s *ExecuteService) resolveBasicHalf(
	field string,
	mapped map[string]string,
	half specType.BasicHalf,
) (string, *config.CustomError) {
	if !strings.Contains(field, "{{secret:") {
		return field, nil
	}
	name, cerr := parseSecretRef(field, mapped)
	if cerr != nil {
		return "", cerr
	}
	sec, cerr := s.SecretService.ResolveSecret(name)
	if cerr != nil {
		return "", cerr
	}
	v, u, p := sec.Reveal()
	if sec.Type() == "basic" {
		if half == specType.BasicHalfUser {
			return u, nil
		}
		return p, nil
	}
	return v, nil
}

// parseSecretRef extracts NAME from "{{secret:NAME}}" and applies the panel
// secretRefs override.
func parseSecretRef(template string, mapped map[string]string) (string, *config.CustomError) {
	m := config.GetSecretPlaceholderRegex().FindStringSubmatch(template)
	if len(m) != 2 {
		return "", config.NewCustomError(
			fmt.Errorf("auth field is not a {{secret:NAME}} placeholder: %q", template),
			config.GetErrCodeInvalidRequest(),
		)
	}
	name := m[1]
	if mapped != nil {
		if v, ok := mapped[name]; ok && v != "" {
			name = v
		}
	}
	return name, nil
}

func headerOrDefault(want, dflt string) string {
	if want == "" {
		return dflt
	}
	return want
}
