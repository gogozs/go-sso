package api_error

var (
	ErrInvalid      = &InvalidError{BaseError{INVALID_PARAMS, INVALID_PARAMS_MSG}}
	ErrUnauthorized = &UnauthorizedError{BaseError{UNAUTHORIZED_ERROR, UNAUTHORIZED_MSG}}
	ErrPermission   = &PermissionError{BaseError{PERMISSION_ERROR, PERMISSION_ERROR_MSG}}
	ErrNotFound     = &NotFoundError{BaseError{NOT_FOUND, NOT_FOUND_MSG}}
	ErrInternal     = &InternalError{BaseError{INTERNAL_ERROR, INTERNAL_ERROR_MSG}}
	// auth
	ErrAuth         = &AuthError{BaseError{AUTH_ERROR, AUTH_ERROR_MSG}}
	ErrTokenExpired = &AuthError{BaseError{TokenExpired_Code, TokenExpired_Msg}}
	ErrTokenNotValidYet = &AuthError{BaseError{TokenNotValidYet_Code, TokenNotValidYet_Msg}}
	ErrTokenMalformed = &AuthError{BaseError{TokenMalformed_Code, TokenMalformed_Msg}}
	ErrTokenInvalid = &AuthError{BaseError{TokenInvalid_Code, TokenInvalid_Msg}}
	ErrSignKey = &AuthError{BaseError{SignKey_Code, SignKey_Msg}}
	// challenge
	ErrInvalidPunch = &InvalidPunch{BaseError{InvalidPunch_Code, InvalidPunch_Msg}}
)

const (
	SUCCESS              = 0
	SUCCESS_MSG          = "请求成功"
	INTERNAL_ERROR       = 500
	INTERNAL_ERROR_MSG   = "Internal Error"
	INVALID_PARAMS       = 400
	INVALID_PARAMS_MSG   = "请求参数错误"
	UNAUTHORIZED_ERROR   = 401
	UNAUTHORIZED_MSG     = "身份认证信息未提供"
	PERMISSION_ERROR     = 403
	PERMISSION_ERROR_MSG = "权限不足，拒绝访问"
	NOT_FOUND            = 404
	NOT_FOUND_MSG        = "请求资源不存在"
	// auth errors
	AUTH_ERROR     = 4001
	AUTH_ERROR_MSG = "账号或密码错误"
	TokenExpired_Code = 4002
	TokenExpired_Msg = "Token is expired"
	TokenNotValidYet_Code = 4003
	TokenNotValidYet_Msg = "Token not active yet"
	TokenMalformed_Code = 4004
	TokenMalformed_Msg = "That's not even a token"
	TokenInvalid_Code = 4005
	TokenInvalid_Msg = "Couldn't handle this token"
	SignKey_Code = 4006
	SignKey_Msg = "newtrekWang"

	// challenge errors
	InvalidPunch_Code = 4101
	InvalidPunch_Msg = "无效的挑战信息"
)
