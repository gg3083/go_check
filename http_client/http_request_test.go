package http_client

import (
	"log"
	"testing"
)

func TestName(t *testing.T) {

	url := "https://i.instagram.com/api/v1/users/formidat_rebe_/usernameinfo/"

	proxySession := ProxySession("http://127.0.0.1:41090")
	httpBack := proxySession.Get(url, nil)
	if httpBack.Code != 0 {
		log.Fatalf("%+v", httpBack)
	}
	data := httpBack.Session.RespData
	log.Println("resp:", string(data))
}
