package api_error

func NewError(err error) ApiError {
	return &BaseError{
		code: 10000,
		msg:  err.Error(),
	}
}

type ApiError interface {
	Error() string
	Code() int
	GetMsg() string
}

type BaseError struct {
	code int
	msg  string
}

func (this *BaseError) Code() int {
	return this.code
}

func (this *BaseError) Error() string {
	return this.msg
}

func (this *BaseError) GetMsg() string {
	return this.Error()
}

type InvalidError struct{ BaseError }
type InternalError struct{ BaseError }
type NotFoundError struct{ BaseError }
type UnauthorizedError struct{ BaseError }
type PermissionError struct{ BaseError }

// auth
type AuthError struct{ BaseError }
type TokenExpired struct{ BaseError }
type TokenNotValidYet struct{ BaseError }
type TokenMalformed struct{ BaseError }
type TokenInvalid struct{ BaseError }
type SignKey struct{ BaseError }

// challenge
type InvalidPunch struct{ BaseError }
