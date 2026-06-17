package generate

type RenderModel struct {
	Method      string
	URL         string
	Headers     []KV
	QueryString string
	BodyEnc     string
	BodyJSON    string
	BodyForm    []KV
	BodyRaw     string
	ContentType string

	AuthType         string
	AuthHeaderName   string
	AuthHeaderValue  string // curl-style: prefix + "$" + envName
	AuthHeaderPrefix string // e.g. "Bearer "
	AuthHeaderEnv    string // e.g. "TRYIT_SECRET_STRIPE_KEY"
	BasicUserEnv     string
	BasicPassEnv     string
	APIKeyIn         string
	APIKeyName       string
	APIKeyEnv        string

	EnvVars []string
}

type KV struct {
	Name  string
	Value string
}
