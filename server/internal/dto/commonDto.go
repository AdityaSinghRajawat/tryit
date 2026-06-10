// Package dto holds the wire shapes for /handler. Handlers do decode→validate
// →one-service-call→render; the dto package is intentionally logic-free.
package dto

// ErrorCode mirrors IMPL §9.2.
type ErrorCode string

const (
	ErrUnauthorized      ErrorCode = "unauthorized"
	ErrForbiddenOrigin   ErrorCode = "forbidden_origin"
	ErrNotPaired         ErrorCode = "not_paired"
	ErrInvalidRequest    ErrorCode = "invalid_request"
	ErrParseFailed       ErrorCode = "parse_failed"
	ErrSecretNotFound    ErrorCode = "secret_not_found"
	ErrAIUnavailable     ErrorCode = "ai_unavailable"
	ErrTargetUnreachable ErrorCode = "target_unreachable"
	ErrInternal          ErrorCode = "internal"
)

type ErrorBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
}

type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}
