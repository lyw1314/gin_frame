package curl

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	//tcp创建超时时长
	connTimeOut time.Duration = 10
	deadTimeOut time.Duration = 10

	//全局HttpTransport,复用TCP连接池
	tr *http.Transport
)

func init() {
	tr = &http.Transport{
		//跳过证书校验https
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, connTimeOut*time.Second)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		IdleConnTimeout:       120 * time.Second,
		MaxConnsPerHost:       2000,
		MaxIdleConns:          5000,
		MaxIdleConnsPerHost:   200,
		DisableKeepAlives:     false,
		ResponseHeaderTimeout: 60 * time.Second,
	}
	//替换默认httpclient的transport配置
	http.DefaultClient.Transport = tr
	http.DefaultClient.Timeout = 60 * time.Second
}

// 发送http请求，method请求方法
// HC=header and Cookie
func httpClient(method, urlHost string, paramStr string, timeout time.Duration, HC ...map[string]string) (string, error) {

	if timeout == 0 {
		timeout = deadTimeOut
	}
	client := &http.Client{
		//Client会在试图用Head、Get、Post或Do方法执行请求时返回错误
		Timeout:   timeout * time.Second,
		Transport: tr,
	}
	var respBodyStr string
	var request *http.Request
	var resp *http.Response
	var err error
	method = strings.ToUpper(method)
	switch method {

	case "GET":
		request, err = http.NewRequest("GET", urlHost, strings.NewReader(paramStr))

	case "POSTQUERY":
		request, err = http.NewRequest("POST", urlHost, strings.NewReader(paramStr))

	case "POSTJSON":
		request, err = http.NewRequest("POST", urlHost, strings.NewReader(paramStr))

	case "PUT":
		request, err = http.NewRequest("PUT", urlHost, strings.NewReader(paramStr))

	}
	for i := 0; i < 1; i++ {
		if err != nil || request == nil {
			break
		}

		if method == "POSTJSON" {
			request.Header.Set("Content-Type", "application/json;charset=UTF-8")
		} else {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:71.0) Gecko/20100101 Firefox/71.0")
		//HC:0:header
		if len(HC) > 0 {
			myHeader := HC[0]
			for hk, hv := range myHeader {
				request.Header.Add(hk, hv)
			}
		}
		resp, err = client.Do(request)
		if err != nil {
			break
		}
		var byteStr []byte
		byteStr, _ = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		//不管是http_status是否是200，仍然使用到返回结果
		//例如mediva、api.crm等
		respBodyStr = string(byteStr)
		//附带响应状态错误信息
		if resp.StatusCode != http.StatusOK {
			//例如504,err默认nil值，需重新改
			err = errors.New("http status:" + strconv.Itoa(resp.StatusCode))
			break
		}
	}
	return respBodyStr, err
}

func Post(post_url string, postData interface{}, timeout time.Duration, HC ...map[string]string) (string, error) {
	var buf bytes.Buffer
	var newData map[string]interface{}
	var method string = "POSTQUERY"
	encoder := json.NewEncoder(&buf)
	encoder.Encode(postData)

	decoder := json.NewDecoder(&buf)
	decoder.UseNumber() // 此处能够保证bigint的精度
	decoder.Decode(&newData)
	values := url.Values{}

	for k, v := range newData {
		if val, ok2 := v.(string); ok2 {
			values.Set(k, val)
		} else if val, ok2 := v.([]interface{}); ok2 {
			for k2, v2 := range val {
				if vv2, ok3 := v2.(string); ok3 {
					kk2 := strconv.Itoa(k2)
					key := k + "[" + kk2 + "]"
					values.Set(key, vv2)
				}
			}
		}
	}
	encodeBody := values.Encode()
	return httpClient(method, post_url, encodeBody, timeout, HC...)
}

func PostJson(post_url string, postData interface{}, timeout time.Duration, HC ...map[string]string) (string, error) {

	var buf bytes.Buffer
	var method string = "POSTJSON"
	encoder := json.NewEncoder(&buf)
	encoder.Encode(postData)
	var jsonBody = buf.String()
	return httpClient(method, post_url, jsonBody, timeout, HC...)
}

// put
func Put(post_url string, post map[string]string, timeout time.Duration, HC ...map[string]string) (string, error) {
	var method string = "PUT"
	values := url.Values{}
	for k, v := range post {
		values.Set(k, v)
	}
	encodeBody := values.Encode()
	return httpClient(method, post_url, encodeBody, timeout, HC...)

}

func Get(get_url string, timeout time.Duration, HC ...map[string]string) (string, error) {
	var method string = "GET"
	return httpClient(method, get_url, "", timeout, HC...)
}

func ValidateUrlType(myUrl string) bool {
	forbidType := []string{"jpg", "jpeg", "gif", "doc", "docx", "exe",
		"mp3", "swf", "txt", "pmb", "avi", "zip", "rar", "wav",
	}
	myUrl = strings.Replace(myUrl, "%", "%25", -1)
	myUrl = strings.ToLower(myUrl)
	u, err := url.Parse(myUrl)
	if err != nil {
		return false
	}
	if u.Path == "" {
		return true
	}
	pathEx := strings.Split(u.Path, ".")
	ext := pathEx[len(pathEx)-1 : len(pathEx)][0]
	if u.RawQuery != "" && ext == "gif" {
		return true
	}
	for _, t := range forbidType {
		if t == ext {
			return false
		}
	}

	return true
}

// 是否是深度链接deeplink
func IsDeepLinkUrl(myUrl string) bool {
	//不允许出现中文、空格、换行符、tab制表符。最多输入1024个字符。
	if len(myUrl) > 1024 {
		return false
	}
	if len(myUrl) != len([]rune(myUrl)) {
		//含有中文,长度不一致了
		return false
	}
	denyChar := []string{" ", "\n", "\t"}
	for _, str := range denyChar {
		if strings.Contains(myUrl, str) {
			return false
		}
	}
	return true
}

func CheckAllowDomain(userInfo map[string]string, aurl string, isMobile ...interface{}) bool {
	aurl = strings.Replace(aurl, "%", "%25", -1)
	//json后的url会有http:\/\/xxx
	aurl = strings.Replace(aurl, "\\", "", -1)
	aurl = strings.ToLower(aurl)
	u, err := url.Parse(aurl)
	if err != nil {
		return false
	}
	host_arr := strings.Split(strings.Split(u.Host, ":")[0], ".")
	if len(host_arr) <= 1 {
		return false
	}
	a := host_arr[len(host_arr)-1]
	b := host_arr[len(host_arr)-2]
	if (a+b) == "cncom" || (a+b) == "auedu" {
		if len(host_arr) == 2 {
			return false
		}
	}

	var domain_list []string
	domain_list = append(domain_list, strings.ToLower(userInfo["website"]))
	domain_list = append(domain_list, strings.Split(strings.ToLower(userInfo["allowwebsite"]), ";")...)
	if len(isMobile) > 0 {
		domain_list = append(domain_list, strings.Split(strings.ToLower(userInfo["mobile_allowwebsite"]), ";")...)
	}

	if len(domain_list) == 0 {
		return false
	}

	for _, v := range domain_list {
		if !strings.HasPrefix(v, "http://") && !strings.HasPrefix(v, "https://") {
			v = "http://" + v
		}
		domain, _ := url.Parse(v)
		if domain.Host == u.Host {
			return true
		}

		host := strings.Split(domain.Host, ".")
		if len(host) <= 1 {
			continue
		}

		c := host[len(host)-1]
		d := host[len(host)-2]

		if (c+d) == "cncom" || (c+d) == "auedu" {
			if len(host) == 2 {
				continue
			}
			if a == c && b == d && host_arr[len(host_arr)-3] == host[len(host)-3] {
				return true
			}
		} else {
			if a == c && b == d {
				return true
			}
		}
	}

	return false
}

// 获取域名url的后缀
// 如果urlList含有非正常url。例如没有.标记的或者.com
func GetUrlSuffix(urlList []string) ([]string, error) {
	var urlSuffix = []string{}
	var err error
	for _, v := range urlList {

		if v == "" {
			continue
		}
		v = strings.Trim(v, " .")
		if v == "" {
			err = errors.New("has err url")
			break
		}

		if !strings.HasPrefix(v, "http://") && !strings.HasPrefix(v, "https://") {
			v = "http://" + v
		}
		domain, _ := url.Parse(v)
		if domain.Host == "" {
			err = errors.New("has err url host")
			break
		}
		host := strings.Split(domain.Host, ".")
		if len(host) <= 1 {
			err = errors.New("has err format url host ")
			break
		}

		sfx := host[len(host)-1]
		if sfx == "" {
			err = errors.New("has err url suffix")
			break
		}
		urlSuffix = append(urlSuffix, sfx)
	}
	return urlSuffix, err
}
