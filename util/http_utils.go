package util

import (
	"crypto/tls"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/coraldane/dnspool"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	CHARSET_UTF8       = "UTF-8"
	DEFAULT_USER_AGENT = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36`
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
	Dial: TimeoutDialer(3*time.Second, 10*time.Second),
	ResponseHeaderTimeout: 30 * time.Second,
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
}
var defaultClient = &http.Client{Transport: defaultTransport}

func DoGetWithRetry(strUrl string, charset string, retryTimes int) (string, error) {
	body, err := GetWithRetry(strUrl, retryTimes)
	if nil != err {
		return "", err
	}
	decoder := mahonia.NewDecoder(charset)
	return decoder.ConvertString(string(body)), nil
}

func GetWithRetry(strUrl string, retryTimes int) (string, error) {
	// fmt.Println("start to get url:", strUrl)
	req, _ := http.NewRequest("GET", strUrl, nil)
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", DEFAULT_USER_AGENT)
	response, err := defaultClient.Do(req)

	if nil != err {
		return getWithRetryLoop(strUrl, retryTimes, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return getWithRetryLoop(strUrl, retryTimes, err)
	}
	// fmt.Println("end to get url:", strUrl)
	strBody := string(body)
	if "null" == strBody || "" == strBody {
		return getWithRetryLoop(strUrl, retryTimes, err)
	}
	return strBody, err
}

func getWithRetryLoop(strUrl string, retryTimes int, err error) (string, error) {
	if retryTimes > 0 {
		return GetWithRetry(strUrl, retryTimes-1)
	} else {
		return "", err
	}
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

func TrimHtmlElement(strSource string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	return re.ReplaceAllString(strSource, "")
}

func TrimBlank(strSource string) string {
	re, _ := regexp.Compile("[\\t\\n]")
	strText := re.ReplaceAllString(strSource, "")
	return strings.Trim(strText, " ")
}

func GetRandomIp() string {
	usefulIpArray := [][]int{
		[]int{607649792, 608174079},     //36.56.0.0-36.63.255.255
		[]int{1038614528, 1039007743},   //61.232.0.0-61.237.255.255
		[]int{1783627776, 1784676351},   //106.80.0.0-106.95.255.255
		[]int{2035023872, 2035154943},   //121.76.0.0-121.77.255.255
		[]int{2078801920, 2079064063},   //123.232.0.0-123.235.255.255
		[]int{-1950089216, -1948778497}, //139.196.0.0-139.215.255.255
		[]int{-1425539072, -1425014785}, //171.8.0.0-171.15.255.255
		[]int{-1236271104, -1235419137}, //182.80.0.0-182.92.255.255
		[]int{-770113536, -768606209},   //210.25.0.0-210.47.255.255
		[]int{-569376768, -564133889},   //222.16.0.0-222.95.255.255
	}

	index := rand.Intn(10)
	intIp := usefulIpArray[index][0] + rand.Intn(usefulIpArray[index][1]-usefulIpArray[index][0])
	return fmt.Sprintf("%d.%d.%d.%d", (intIp>>24)&0xff, (intIp>>16)&0xff, (intIp>>8)&0xff, intIp&0xff)
}
