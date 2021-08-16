package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Method string
type Params map[string]interface{}

type Session struct {
	client    *http.Client
	Header    http.Header
	RespCode  int
	RespData  []byte
	Cookie    []*http.Cookie
	notHeader bool
}

func (session *Session) defaultClient() {
	session.client = http.DefaultClient
}

func DefaultSession() *Session {
	session := Session{}
	session.defaultClient()
	return &session
}

func ProxySession(proxy string) *Session {
	if proxy == "" || !strings.HasPrefix(proxy, "http") {
		return DefaultSession()
	}

	u := url.URL{}
	session := Session{}

	urlProxy, _ := u.Parse(proxy)
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		},
		Timeout: time.Duration(15) * time.Second,
	}
	session.client = &c
	return &session
}

func (session *Session) ClientProxy(proxy string) {
	if proxy == "" || !strings.HasPrefix(proxy, "http") {
		session.defaultClient()
		return
	}

	u := url.URL{}
	urlProxy, _ := u.Parse(proxy)
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		},
		Timeout: time.Duration(15) * time.Second,
	}
	session.client = &c
}

func (session *Session) SetHeader(hdr http.Header) {
	session.Header = hdr
}

func (session *Session) AddHeader(maps map[string]string) {
	header := session.Header
	if header == nil {
		header = http.Header{}
	}
	for k, v := range maps {
		header.Set(k, v)
	}
	session.Header = header
}

func (session *Session) SetCookie(cookie string) {
	hdr := http.Header{}
	if session.Header != nil {
		hdr = session.Header
	}
	hdr.Set("cookie", cookie)
	session.Header = hdr
}

func (session *Session) Get(path string, params Params) HttpBack {
	return session.Api(path, http.MethodGet, params)
}

func (session *Session) Post(path string, params Params) HttpBack {
	return session.Api(path, http.MethodPost, params)
}

func (session *Session) PostForUrl(path string, params Params) HttpBack {
	session.notHeader = true
	return session.Api(path, http.MethodPost, params)
}

func (session *Session) PostForJson(path string, params Params) HttpBack {
	header := session.Header
	if header == nil {
		header = http.Header{}
	}
	header.Set("Content-Type", "application/json")
	session.SetHeader(header)
	return session.Api(path, http.MethodPost, params)
}

func (session *Session) Api(path string, method Method, params Params) HttpBack {
	return session.graph(path, method, params)
}

func (session *Session) graph(path string, method Method, params Params) HttpBack {

	if params == nil {
		params = Params{}
	}
	if method == http.MethodGet {
		path = fmt.Sprintf("%s%s", path, urlAppendParam(params))
		return session.sendGetRequest(path)

	} else if method == http.MethodPost {
		return session.sendPostRequest(path, params)

	}

	return ErrorBack(ERROR)
}

func (session *Session) sendGetRequest(uri string) HttpBack {
	//log.Printf("请求的接口为 %s\n", uri)
	//log.Printf("请求头为 %v\n", session.Header)
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return ErrorBack(URL_ERROR)
	}
	req := &http.Request{
		Method:     http.MethodGet,
		URL:        parsedURL,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     session.Header,
	}

	response, data, httpBack := session.sendRequest(req)

	if httpBack.Code != SUCCESS {
		return httpBack
	}
	session.RespData = data
	session.Cookie = response.Cookies()
	return SuccessSessionBack(session)
}
func (session *Session) sendPostRequest(uri string, params Params) HttpBack {

	var rc io.Reader
	//log.Printf("请求的接口为 %s\n", uri)
	if session.Header == nil {
		session.Header = http.Header{}
	}
	contentType := session.Header.Get("Content-Type")

	if session.notHeader {
		uri = fmt.Sprintf("%s%s", uri, urlAppendParam(params))
	} else if strings.Contains(contentType, "json") {
		jsonParams, err := json.Marshal(params)
		if err != nil {
			err2 := fmt.Errorf("post params json encode error： %v", err)
			return ErrorMsgBack(NETWORK_REQ_PARAM_ERROR, err2.Error())
		}
		rc = bytes.NewReader(jsonParams)
	} else {
		if contentType == "" {
			contentType = "application/x-www-form-urlencoded"
			session.Header.Set("Content-Type", contentType)
		}
		rc = createFormReader(params)
	}
	request, err := http.NewRequest(http.MethodPost, uri, rc)
	if err != nil {
		return ErrorMsgBack(NETWORK_REQ_ERROR, err.Error())
	}
	request.Header = session.Header

	//marshal, _ := json.Marshal(session.Header)
	//log.Printf("请求类型为 %v\n", contentType)
	//log.Printf("请求头为 %v\n", string(marshal))

	response, data, httpBack := session.sendRequest(request)

	if httpBack.Code != SUCCESS {
		return httpBack
	}
	session.RespData = data
	session.Cookie = response.Cookies()
	session.RespCode = response.StatusCode

	//cookie, _ := json.Marshal(response.Cookies())
	//log.Println("返回cookie:", string(cookie))
	log.Println("====================")
	if response.StatusCode != 200 {
		return ErrorMsgBack(NETWORK_RESP_STATUS_ERRPR, fmt.Sprintf("原始状态码：%d", response.StatusCode))
	}
	return SuccessSessionBack(session)
}

func (session *Session) sendRequest(request *http.Request) (*http.Response, []byte, HttpBack) {

	var err error
	var response *http.Response
	var data []byte

	request.Close = false
	startTime := time.Now().UnixNano()
	if session.client == nil {
		response, err = http.DefaultClient.Do(request)
	} else {
		response, err = session.client.Do(request)
	}
	log.Printf("http！耗时:%dms\n", (time.Now().UnixNano()-startTime)/1e6)

	if err != nil {
		if strings.Contains(err.Error(), "proxyconnect") {
			return nil, nil, ErrorMsgBack(PROXY_URL_ERROR, err.Error())
		}
		err = fmt.Errorf("网络异常,发送请求失败: %v", err)
		return nil, nil, ErrorMsgBack(NETWORK_ERROR, err.Error())
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		err = fmt.Errorf("copy result error : %v", err)
		return nil, nil, ErrorMsgBack(NETWORK_RESP_RESULT_ERRPR, err.Error())

	}
	err = response.Body.Close()

	if err != nil {
		err = fmt.Errorf("close http client : %v", err)
		return nil, nil, ErrorMsgBack(NETWORK_RESP_CLOSE_ERRPR, err.Error())
	}

	data = buf.Bytes()
	//log.Printf("[返回]: code:%s , %s\n", response.Status, string(data))
	return response, data, SuccessBack()
}
