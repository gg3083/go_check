package main

import (
	"go_check_proxy/ui"
	"log"
	"net"
)

//const reqUrl = "https://www.baidu.com"
//go:generate go run -tags generate gen.go
func main() {
	//
	//fileName := "index"
	//
	//port, _ := PickUnusedPort()
	//ui.ShowUI(port, fileName)
	// 4780
	//startProxy := ui.StartProxy()
	//log.Printf("%+v\n", startProxy)

	proxy := ui.ReqForProxy("127.0.0.1:4780")
	log.Printf("%+v\n", proxy)

}

func PickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}
