package httplib

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"github.com/axgle/mahonia"
	"github.com/coraldane/dnspool"
)

const (
	CHARSET_UTF8       = "UTF-8"
	DEFAULT_USER_AGENT = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36`
)

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		//conn, err := net.DialTimeout(netw, addr, cTimeout)
		conn, err := dnspool.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

var defaultTransport = &http.Transport{
	Dial:                  TimeoutDialer(3*time.Second, 10*time.Second),
	ResponseHeaderTimeout: 30 * time.Second,
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
}
var defaultClient = &http.Client{Transport: defaultTransport}

func GetWithRetryAndDecode(strUrl string, charset string, retryTimes int) (strResponse string, err error) {
	var (
		body []byte
	)
	body, err = GetWithRetry(strUrl, retryTimes)
	if nil != err {
		return
	}
	decoder := mahonia.NewDecoder(charset)
	strResponse = decoder.ConvertString(string(body))
	return
}

func GetWithRetry(strUrl string, retryTimes int) (body []byte, err error) {
	if 0 > retryTimes {
		err = errors.New("exceed max retry times")
		return
	}

	body, err = Get(strUrl)
	if nil != err {
		goto ERR
	}

ERR:
	GetWithRetry(strUrl, retryTimes-1)
	return
}

func Get(strUrl string) (body []byte, err error) {
	var (
		response *http.Response
	)
	req, _ := http.NewRequest("GET", strUrl, nil)
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", DEFAULT_USER_AGENT)
	response, err = defaultClient.Do(req)
	if nil != err {
		return
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	return
}

func PostForm(method, strUrl string, data url.Values, headers map[string]string) ([]byte, error) {
	if nil == headers {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return HttpRequest(method, strUrl, strings.NewReader(data.Encode()), headers)
}

func HttpRequest(method, strUrl string, data io.Reader, headers map[string]string) ([]byte, error) {
	req, _ := http.NewRequest(strings.ToUpper(method), strUrl, data)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", DEFAULT_USER_AGENT)
	if nil != headers {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
	}
	resp, err := defaultClient.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func DoHttpRequest(method, strUrl string, data io.Reader, charset string, headers map[string]string) (string, error) {
	body, err := HttpRequest(method, strUrl, data, headers)
	if nil != err {
		return "", err
	}

	strBody := string(body)
	if "" != charset {
		decoder := mahonia.NewDecoder(charset)
		strBody = decoder.ConvertString(strBody)
	}
	return strBody, nil
}
