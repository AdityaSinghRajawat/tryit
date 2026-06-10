package execute

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/AdityaSinghRajawat/tryit/server/internal/helpers"
	"github.com/AdityaSinghRajawat/tryit/server/internal/model"
)

// SecretResolver is the port the auth injector uses to obtain a Secret value
// by NAME at the moment of sending. service/secret.Service satisfies it.
type SecretResolver interface {
	Resolve(name string) (model.Secret, error)
}

var placeholderRe = regexp.MustCompile(`\{\{secret:([A-Z0-9_]+)\}\}`)

type basicHalf int

const (
	basicHalfUser basicHalf = iota
	basicHalfPass
)

// injectAuth resolves placeholders and stamps the outbound request. mappedRefs
// maps placeholder NAME → stored secret name (the panel's secretRefs). When a
// placeholder has no entry we resolve by its own name (identity).
//
// Returns a masked preview string for the request preview (never the real
// value).
func injectAuth(req *http.Request, spec model.RequestSpec, mappedRefs map[string]string, r SecretResolver) (maskedAuthPreview string, err error) {
	a := spec.Auth
	switch model.AuthType(a.Type) {
	case model.AuthNone, "":
		return "", nil

	case model.AuthBearer:
		val, err := resolveScalar(a.ValueRef, mappedRefs, r)
		if err != nil {
			return "", err
		}
		header := headerName(a.Name, "Authorization")
		prefix := a.Prefix
		if prefix == "" {
			prefix = "Bearer "
		}
		req.Header.Set(header, prefix+val)
		return header + ": " + helpers.MaskBearer(prefix+val), nil

	case model.AuthAPIKey:
		val, err := resolveScalar(a.ValueRef, mappedRefs, r)
		if err != nil {
			return "", err
		}
		header := headerName(a.Name, "")
		if header == "" {
			return "", errors.New("apiKey auth requires a name")
		}
		switch a.In {
		case "query":
			u := req.URL
			q := u.Query()
			q.Set(header, a.Prefix+val)
			u.RawQuery = q.Encode()
			return "?" + header + "=" + helpers.Mask(val), nil
		case "header", "":
			req.Header.Set(header, a.Prefix+val)
			return header + ": " + helpers.Mask(val), nil
		default:
			return "", fmt.Errorf("invalid apiKey auth.in %q", a.In)
		}

	case model.AuthBasic:
		user, err := resolveBasicHalf(a.Username, mappedRefs, r, basicHalfUser)
		if err != nil {
			return "", err
		}
		pass, err := resolveBasicHalf(a.Password, mappedRefs, r, basicHalfPass)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Basic "+helpers.BasicAuthValue(user, pass))
		return "Authorization: Basic ••••", nil

	default:
		return "", fmt.Errorf("unknown auth.type %q", a.Type)
	}
}

// resolveScalar resolves a bearer/apiKey valueRef ("{{secret:NAME}}") to the
// underlying string. The Secret's underlying type is whatever the store
// returned (typically "bearer" for env-backed Phase 1).
func resolveScalar(template string, mapped map[string]string, r SecretResolver) (string, error) {
	name, err := refName(template, mapped)
	if err != nil {
		return "", err
	}
	sec, err := r.Resolve(name)
	if err != nil {
		return "", err
	}
	v, _, _ := sec.Reveal()
	return v, nil
}

// resolveBasicHalf handles either "{{secret:NAME}}" placeholders or literal
// values. When the placeholder resolves to a basic-typed Secret, we return
// the user/pass half explicitly. When it resolves to a scalar Secret, we
// return its value (Twilio-style: two distinct stored secrets).
func resolveBasicHalf(field string, mapped map[string]string, r SecretResolver, half basicHalf) (string, error) {
	if !strings.Contains(field, "{{secret:") {
		return field, nil
	}
	name, err := refName(field, mapped)
	if err != nil {
		return "", err
	}
	sec, err := r.Resolve(name)
	if err != nil {
		return "", err
	}
	v, u, p := sec.Reveal()
	if sec.Type() == "basic" {
		if half == basicHalfUser {
			return u, nil
		}
		return p, nil
	}
	return v, nil
}

func refName(template string, mapped map[string]string) (string, error) {
	m := placeholderRe.FindStringSubmatch(template)
	if len(m) != 2 {
		return "", fmt.Errorf("auth field is not a {{secret:NAME}} placeholder: %q", template)
	}
	name := m[1]
	if mapped != nil {
		if v, ok := mapped[name]; ok && v != "" {
			name = v
		}
	}
	return name, nil
}

func headerName(want, dflt string) string {
	if want == "" {
		return dflt
	}
	return want
}
