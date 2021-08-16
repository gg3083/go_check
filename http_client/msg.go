package http_client

var ResultMap = map[int]string{
	SUCCESS:         "操作成功！",
	ERROR:           "操作失败！",
	URL_ERROR:       "请求的地址有误！",
	PROXY_URL_ERROR: "设置的代理有误！",

	NETWORK_ERROR:           "网络错误！",
	NETWORK_REQ_ERROR:       "请求错误！",
	NETWORK_REQ_PARAM_ERROR: "请求参数异常！",

	NETWORK_RESP_ERRPR:        "请求返回异常！",
	NETWORK_RESP_RESULT_ERRPR: "请求结果解析异常！",
	NETWORK_RESP_CLOSE_ERRPR:  "请求关闭异常！",
	NETWORK_RESP_STATUS_ERRPR: "请求状态码异常！",
}

func GetMsg(code int) string {
	msg, ok := ResultMap[code]
	if ok {
		return msg
	}

	return ResultMap[ERROR]
}
