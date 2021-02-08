package apierror

func WrapError(err error) ApiError {
	return &BaseError{
		code: 10000,
		msg:  err.Error(),
	}
}

func NewParamsError(msg string) ApiError {
	return &BaseError{
		code: 400,
		msg:  msg,
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
