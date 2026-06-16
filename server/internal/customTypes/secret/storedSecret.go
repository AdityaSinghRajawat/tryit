package secret

// StoredSecret is the persistence record shared between the secret service
// and its store implementations. Value is set for bearer/apiKey; User+Pass
// for basic. Never serialised onto the wire — only the masking Secret and
// Info types are.
type StoredSecret struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	HostHint string `json:"hostHint,omitempty"`
	Value    string `json:"value,omitempty"`
	User     string `json:"user,omitempty"`
	Pass     string `json:"pass,omitempty"`
}
