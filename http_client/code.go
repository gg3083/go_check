package http_client

const (
	SUCCESS = 0
	ERROR   = 1

	URL_ERROR       = 10000
	PROXY_URL_ERROR = 10001

	NETWORK_ERROR           = 11000
	NETWORK_REQ_ERROR       = 12001
	NETWORK_REQ_PARAM_ERROR = 12002

	NETWORK_RESP_ERRPR        = 13001
	NETWORK_RESP_RESULT_ERRPR = 13002
	NETWORK_RESP_CLOSE_ERRPR  = 13003
	NETWORK_RESP_STATUS_ERRPR = 13004
)

type HttpBack struct {
	Code        int    `json:"code"`
	BizMsg      string `json:"biz_msg"`
	OriginalMsg string `json:"original_msg"`
	*Session
}

func SuccessBack() HttpBack {
	return HttpBack{
		Code:   SUCCESS,
		BizMsg: GetMsg(SUCCESS),
	}
}

func SuccessSessionBack(session *Session) HttpBack {
	return HttpBack{
		Code:    SUCCESS,
		BizMsg:  GetMsg(SUCCESS),
		Session: session,
	}
}

func ErrorBack(code int) HttpBack {
	return HttpBack{
		Code:   code,
		BizMsg: GetMsg(code),
	}
}
func ErrorMsgBack(code int, msg string) HttpBack {
	return HttpBack{
		Code:        code,
		BizMsg:      GetMsg(code),
		OriginalMsg: msg,
	}
}
