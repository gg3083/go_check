package ui

import "encoding/json"

type JsonBack struct {
	Code   int         `json:"code"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

func Success() string {
	return SuccessData(nil)
}

func SuccessData(data interface{}) string {
	back := JsonBack{
		Code: 0,
		Data: data,
	}
	res, _ := json.Marshal(back)
	return string(res)
}

func SuccessDataMsg(data interface{}, msg string) string {
	back := JsonBack{
		Code:   0,
		Data:   data,
		ErrMsg: msg,
	}
	res, _ := json.Marshal(back)
	return string(res)
}

func Error() string {
	return ErrorMsg("系统异常")
}
func ErrorMsg(errMsg string) string {
	back := JsonBack{
		Code:   1,
		ErrMsg: errMsg,
	}
	res, _ := json.Marshal(back)
	return string(res)
}
