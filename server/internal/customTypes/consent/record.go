package consent

import "time"

// Record: one persisted entry per (storedSecretName, targetHost) granted.
type Record struct {
	Secret    string    `json:"secret"`
	Host      string    `json:"host"`
	GrantedAt time.Time `json:"grantedAt"`
}
