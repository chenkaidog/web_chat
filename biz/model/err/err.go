package err

import "fmt"

type Error interface {
	Error() string
	Code() int32
	Msg() string
	SetErr(err error) Error
	SetMsg(msg string) Error
}

type bizError struct {
	code int32
	msg  string
}

func (err *bizError) Error() string {
	return fmt.Sprintf("%d:%s", err.code, err.msg)
}

func (err *bizError) Code() int32 {
	return err.code
}

func (err *bizError) Msg() string {
	return err.msg
}

func (bizErr *bizError) SetErr(err error) Error {
	return New(bizErr.Code(), err.Error())
}

func (bizErr *bizError) SetMsg(msg string) Error {
	return New(bizErr.Code(), msg)
}

func New(code int32, msg string) Error {
	return &bizError{
		code: code,
		msg:  msg,
	}
}

func ErrorEqual(err1, err2 Error) bool {
	// 都为空
	if err1 == nil && err2 == nil {
		return true
	}

	// 只有一个不为空
	if err1 == nil || err2 == nil {
		return false
	}

	// 都不为空
	return err1.Code() == err2.Code()
}

var (
	Success             = New(0, "success")
	ServerError         = New(1_0001, "internal server error")
	ParamError          = New(1_0002, "param error")
	InternalServerError = New(1_000_3, "internal server error")

	AccountNotExistError      = New(2_0001, "account not exist")
	PasswordIncorrect         = AccountNotExistError
	AccountStatusInvalidError = New(2_0002, "account is invalid")
	AccountExpiredError       = New(2_000_3, "account is expired")
)
