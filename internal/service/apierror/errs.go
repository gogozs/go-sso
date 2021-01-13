package apierror

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

func (e *BaseError) Code() int {
	return e.code
}

func (e *BaseError) Error() string {
	return e.msg
}

func (e *BaseError) GetMsg() string {
	return e.Error()
}
