package ui

import (
	"fmt"
	"go_check_proxy/http_client"
)

//const ReqUrl = "https://i.instagram.com/api/v1/users/formidat_rebe_/usernameinfo/"
const ReqUrl2 = "https://google.com"

func ReqForProxy(proxyUrl string) string {
	if proxyUrl == "" {
		return ErrorMsg("代理地址不能为空")
	}
	proxyUrl = fmt.Sprintf("http://%s", proxyUrl)
	session := http_client.ProxySession(proxyUrl)
	httpBack := session.Get(ReqUrl2, nil)
	if httpBack.Code != 0 {
		return ErrorMsg(httpBack.BizMsg)
	}
	return SuccessData(string(httpBack.RespData))
}
