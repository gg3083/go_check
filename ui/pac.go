package ui

import (
	"fmt"
	"go_check_proxy/http_client"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"regexp"
)

func StartProxy() string {
	proxyUrl, pacUrl, err := Regedit()
	if err != nil {
		return ErrorMsg("未检测到代理：" + err.Error())
	}
	if proxyUrl != "" {
		return SuccessData(proxyUrl)
	}
	res := ReadPacUrl(pacUrl)
	if res != "" {
		return SuccessData(res)
	}
	return ErrorMsg("读取代理失败")
}

func ReadPacFile(fileName string) string {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("读取文件失败:", err.Error())
		return ""
	}

	return ParsePac(file)
}

func ReadPacUrl(url string) string {
	session := http_client.DefaultSession()
	httpBack := session.Get(url, nil)
	if httpBack.Code != http_client.SUCCESS {
		return ""
	}
	log.Printf("res:%+v\n", httpBack)
	return ParsePac(httpBack.Session.RespData)
}

func ParsePac(content []byte) string {
	reg := "127.0.0.1:([0-9]|[1-9]\\d{1,3}|[1-5]\\d{4}|6[0-4]\\d{4}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])*"
	exp := regexp.MustCompile(reg)
	allString := exp.FindAllString(string(content), -1)
	if len(allString) > 0 {
		return allString[0]
	}
	return ""
}

func Regedit() (string, string, error) {
	//HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings
	// ProxyEnable
	// AutoConfigURL
	openKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		return "", "", fmt.Errorf("读取注册表失败！%s", err.Error())
	}
	proxyEnableValue, valtype1, err := openKey.GetIntegerValue("ProxyEnable")
	if err != nil {
		return "", "", err
	}
	log.Printf("proxyEnableValue:%v, valtype:%v,\n", proxyEnableValue, valtype1)

	//开启了
	if proxyEnableValue == 1 {
		proxyServerValue, _, err := openKey.GetStringValue("ProxyServer")
		if err != nil {
			return "", "", err
		}
		return proxyServerValue, "", nil
	}

	stringsValue, valtype, err := openKey.GetStringValue("AutoConfigURL")
	if err != nil {
		return "", "", fmt.Errorf("读取AutoConfigURL异常：%s", err.Error())
	}
	log.Printf("stringsValue:%v, valtype:%v,\n", stringsValue, valtype)
	// ProxyServer ProxyEnable
	return "", stringsValue, nil

}
