package secret

// Info: names + metadata only, never values.
type Info struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	HostHint string `json:"hostHint,omitempty"`
}

type ListResponse struct {
	Secrets []Info `json:"secrets"`
}
