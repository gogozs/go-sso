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
