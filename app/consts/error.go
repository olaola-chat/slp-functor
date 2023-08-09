package consts

type CommonError interface {
	Code() int32
	Msg() string
	Error() string
}

type commonError struct {
	code int32
	msg  string
}

func (e *commonError) Code() int32 {
	return e.code
}

func (e *commonError) Msg() string {
	return e.msg
}

func (e *commonError) Error() string {
	return e.msg
}

func newError(code int32, msg string) CommonError {
	e := &commonError{
		code: code,
		msg:  msg,
	}
	return e
}

var (
	ERROR_SUCCESS = newError(0, "") // 成功

	// 10001-10099 基本错误或系统级错误
	ERROR_SYSTEM = newError(10000, "服务器异常，请稍后重试") // 系统未知错误
	ERROR_PARAM  = newError(10002, "参数异常")        // 参数解析错误
)
