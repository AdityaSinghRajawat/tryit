package config

type CustomErrorCode string

// CustomError carries an inner error plus a classification code; does not
// implement the error interface so callers can't accidentally drop ErrCode.
type CustomError struct {
	Error   error           `json:"error"`
	ErrCode CustomErrorCode `json:"errCode"`
}

func (e *CustomError) Is(target *CustomError) bool {
	if e == nil || target == nil {
		return e == target
	}
	return e.Error == target.Error && e.ErrCode == target.ErrCode
}

func NewCustomError(err error, errCode CustomErrorCode) *CustomError {
	return &CustomError{Error: err, ErrCode: errCode}
}

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

func GetErrCodeUnauthorized() CustomErrorCode      { return errCodesI.Shared.errCodeUnauthorized }
func GetErrCodeForbiddenOrigin() CustomErrorCode   { return errCodesI.Shared.errCodeForbiddenOrigin }
func GetErrCodeNotPaired() CustomErrorCode         { return errCodesI.Shared.errCodeNotPaired }
func GetErrCodeInvalidRequest() CustomErrorCode    { return errCodesI.Shared.errCodeInvalidRequest }
func GetErrCodeParseFailed() CustomErrorCode       { return errCodesI.Shared.errCodeParseFailed }
func GetErrCodeSecretNotFound() CustomErrorCode    { return errCodesI.Shared.errCodeSecretNotFound }
func GetErrCodeAIUnavailable() CustomErrorCode     { return errCodesI.Shared.errCodeAIUnavailable }
func GetErrCodeTargetUnreachable() CustomErrorCode { return errCodesI.Shared.errCodeTargetUnreachable }
func GetErrCodeInternal() CustomErrorCode          { return errCodesI.Shared.errCodeInternal }
