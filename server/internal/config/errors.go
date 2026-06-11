package config

// CustomErrorCode is the typed string carried through every service/handler
// to classify an error. Values are SCREAMING_SNAKE_CASE so wire payloads stay
// consistent across services.
type CustomErrorCode string

// CustomError is the standard error carrier returned by services and consumed
// by handlers via utils.HandleCustomError. It deliberately does NOT implement
// the `error` interface — callers always pass *CustomError explicitly so the
// ErrCode classification is never lost.
type CustomError struct {
	Error   error           `json:"error"`
	ErrCode CustomErrorCode `json:"errCode"`
}

// Is compares two *CustomError values structurally. Not the same as
// errors.Is — it takes a *CustomError, not error, and is used at call sites
// where both sides are known to be CustomErrors.
func (e *CustomError) Is(target *CustomError) bool {
	if e == nil || target == nil {
		return e == target
	}
	return e.Error == target.Error && e.ErrCode == target.ErrCode
}

func NewCustomError(err error, errCode CustomErrorCode) *CustomError {
	return &CustomError{Error: err, ErrCode: errCode}
}

// genericErrCodes holds codes used across all domains. Future domain-specific
// groups slot in next to it as additional anonymous fields on `errCodes`
// (e.g. `Pairing`, `Execution`).
type genericErrCodes struct {
	errCodeUnauthorized      CustomErrorCode
	errCodeForbiddenOrigin   CustomErrorCode
	errCodeNotPaired         CustomErrorCode
	errCodeInvalidRequest    CustomErrorCode
	errCodeParseFailed       CustomErrorCode
	errCodeSecretNotFound    CustomErrorCode
	errCodeAIUnavailable     CustomErrorCode
	errCodeTargetUnreachable CustomErrorCode
	errCodeInternal          CustomErrorCode
}

type errCodes struct {
	Shared genericErrCodes
}

var errCodesI = &errCodes{
	Shared: genericErrCodes{
		errCodeUnauthorized:      "UNAUTHORIZED",
		errCodeForbiddenOrigin:   "FORBIDDEN_ORIGIN",
		errCodeNotPaired:         "NOT_PAIRED",
		errCodeInvalidRequest:    "INVALID_REQUEST",
		errCodeParseFailed:       "PARSE_FAILED",
		errCodeSecretNotFound:    "SECRET_NOT_FOUND",
		errCodeAIUnavailable:     "AI_UNAVAILABLE",
		errCodeTargetUnreachable: "TARGET_UNREACHABLE",
		errCodeInternal:          "INTERNAL_SERVER_ERROR",
	},
}

func GetErrCodeUnauthorized() CustomErrorCode { return errCodesI.Shared.errCodeUnauthorized }

func GetErrCodeForbiddenOrigin() CustomErrorCode { return errCodesI.Shared.errCodeForbiddenOrigin }
func GetErrCodeNotPaired() CustomErrorCode       { return errCodesI.Shared.errCodeNotPaired }
func GetErrCodeInvalidRequest() CustomErrorCode  { return errCodesI.Shared.errCodeInvalidRequest }
func GetErrCodeParseFailed() CustomErrorCode     { return errCodesI.Shared.errCodeParseFailed }
func GetErrCodeSecretNotFound() CustomErrorCode  { return errCodesI.Shared.errCodeSecretNotFound }
func GetErrCodeAIUnavailable() CustomErrorCode   { return errCodesI.Shared.errCodeAIUnavailable }

func GetErrCodeTargetUnreachable() CustomErrorCode { return errCodesI.Shared.errCodeTargetUnreachable }
func GetErrCodeInternal() CustomErrorCode          { return errCodesI.Shared.errCodeInternal }
