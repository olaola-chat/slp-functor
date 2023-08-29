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

	ERROR_AUDIT_AUDIO_PARAM            = newError(10003, "审核音频状态参数错误")
	ERROR_AUDIO_NOT_EXIST              = newError(10004, "音频不存在")
	ERROR_AUDIO_COLLECT                = newError(10005, "该音频处于收录状态")
	ERROR_ALBUM_HAS_AUDIO              = newError(10006, "该专辑已收录音频")
	ERROR_ALBUM_NOT_EXIST              = newError(10007, "专辑不存在")
	ERROR_AUDIO_STATUS_INVALID         = newError(10008, "该音频审核状态未通过")
	ERROR_AUDIO_ALBUM_COLLECT          = newError(10009, "该音频已收录在该专辑下")
	ERROR_AUDIO_ALBUM_COLLECT_REMOVE   = newError(10009, "该音频还未收录在该专辑下")
	ERROR_ALBUM_COLLECT                = newError(10010, "该专辑处于收录状态")
	ERROR_SUBJECT_HAS_ALBUM            = newError(10011, "该主题已收录专辑")
	ERROR_ALBUM_SUBJECT_COLLECT        = newError(10012, "该专辑已收录在该专题下")
	ERROR_ALBUM_SUBJECT_COLLECT_REMOVE = newError(10013, "该专辑还未收录在该专题下")
	ERROR_SUBJECT_NOT_EXIST            = newError(10014, "专题不存在")
)
