package http_client

import (
	"fmt"
	"io"
	"net/url"
	"strings"
)

// 构建form表单参数
func createFormReader(params Params) io.Reader {
	form := url.Values{}
	for k, v := range params {
		form.Add(k, fmt.Sprintf("%v", v))
	}
	return strings.NewReader(form.Encode())
}

// 构建Url拼接参数
func urlAppendParam(param map[string]interface{}) string {
	if param == nil {
		return ""
	}
	urlParams := ""
	for k, v := range param {
		if v != "" {
			urlParams = fmt.Sprintf("%v%v=%v&", urlParams, k, v)
		}
	}
	if urlParams != "" {
		return "?" + urlParams[0:len(urlParams)-1]
	}
	return urlParams
}
