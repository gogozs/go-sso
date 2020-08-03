package api_error

var (
	ErrInvalid          = &BaseError{InvalidParamsCode, InvalidParamsMsg}
	ErrUnauthorized     = &BaseError{UnauthorizedErrorCode, UnauthorizedMsg}
	ErrPermission       = &BaseError{PermissionErrorCode, PermissionErrorMsg}
	ErrNotFound         = &BaseError{NotFoundCode, NotFoundMsg}
	ErrInternal         = &BaseError{InternalErrorCode, InternalErrorMsg}
	ErrAuth             = &BaseError{AuthErrorCode, AuthErrorMsg}
	ErrTokenInvalid     = &BaseError{TokenInvalidCode, TokenInvalidMsg}
	ErrTokenExpired     = &BaseError{TokenExpiredCode, TokenExpiredMsg}
	ErrTokenNotValidYet = &BaseError{TokenNotValidYetCode, TokenNotValidYetMsg}
	ErrTokenMalformed   = &BaseError{TokenMalformedCode, TokenMalformedMsg}
)

const (
	// common
	SuccessCode           = 0
	SuccessMsg            = "请求成功"
	InternalErrorCode     = 500
	InternalErrorMsg      = "Internal Error"
	InvalidParamsCode     = 400
	InvalidParamsMsg      = "请求参数错误"
	UnauthorizedErrorCode = 401
	UnauthorizedMsg       = "身份认证信息未提供"
	PermissionErrorCode   = 403
	PermissionErrorMsg    = "权限不足，拒绝访问"
	NotFoundCode          = 404
	NotFoundMsg           = "请求资源不存在"
	AuthErrorCode         = 4000
	AuthErrorMsg          = "auth error"
	TokenInvalidCode      = 4001
	TokenInvalidMsg       = "invalid token"
	TokenExpiredCode      = 4002
	TokenExpiredMsg       = "token expired"
	TokenNotValidYetCode  = 4003
	TokenNotValidYetMsg   = "token not valid yet"
	TokenMalformedCode    = 4004
	TokenMalformedMsg     = "That's not even a token"
)
