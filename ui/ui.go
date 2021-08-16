package ui

import (
	"fmt"
	"github.com/zserge/lorca"
	assets "go_check_proxy/assets_gen"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
)

var LocalUi lorca.UI
var err error

func ShowUI(port int, fileName string) {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	LocalUi, err = lorca.New("", "", 960, 480, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer LocalUi.Close()

	LocalUi.Bind("goFuncLoadProxy", StartProxy)
	LocalUi.Bind("goFuncPingProxy", ReqForProxy)

	LocalUi.Bind("goFuncExit", func() {
		log.Println("系统退出！！！")
		LocalUi.Close()
		os.Exit(1)
	})

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(assets.FS))

	addr := fmt.Sprintf("http://%s/%s.html", ln.Addr(), fileName)
	log.Println(addr)
	LocalUi.Load(addr)

	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	LocalUi.Eval(`
		console.log("当前版本v1.0.2");
	`)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-LocalUi.Done():
	}

	log.Println("exiting...")

}
