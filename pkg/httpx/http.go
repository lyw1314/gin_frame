package httpx

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
)

var (
	locker = sync.RWMutex{}
	pool   = make(map[Config]*http.Client)
)

func init() {
	locker.Lock()
	pool[defaultConfig] = &http.Client{Transport: Transport(defaultConfig)}
	locker.Unlock()
}

func PostForm(ctx context.Context, url string, data url.Values, headers ...Header) (resp []byte, err error) {
	return PostFormX(ctx, defaultConfig, url, data, headers...)
}

func PostJson(ctx context.Context, url string, data []byte, headers ...Header) (resp []byte, err error) {
	return PostJsonX(ctx, defaultConfig, url, data, headers...)
}
func Post(ctx context.Context, url string, data []byte, contentType string, headers ...Header) (resp []byte, err error) {
	return PostX(ctx, defaultConfig, url, data, contentType, headers...)
}

// 支持协程并发
func Get(ctx context.Context, url string, headers ...Header) (resp []byte, err error) {
	return GetX(ctx, defaultConfig, url, headers...)
}

func BuildQuery(api string, params url.Values) string {
	if params != nil {
		Url, _ := url.ParseRequestURI(api)
		p := Url.Query()
		for k, values := range params {
			for _, v := range values {
				p.Add(k, v)
			}
		}
		Url.RawQuery = p.Encode()
		return Url.String()
	}
	return api
}

func tryAgain(times int, caller func() (resp []byte, err error)) (resp []byte, err error) {
	for i := 0; i < times-1; i++ {
		if resp, err := caller(); err == nil {
			return resp, err
		}
	}
	return caller()
}

func PostFormX(ctx context.Context, config Config, url string, data url.Values, headers ...Header) (resp []byte, err error) {
	return commonPost(ctx, getClient(config), url, []byte(data.Encode()), "application/x-www-form-urlencoded", headers...)
}

func PostJsonX(ctx context.Context, config Config, url string, data []byte, headers ...Header) (resp []byte, err error) {
	return commonPost(ctx, getClient(config), url, data, "application/json", headers...)
}

func GetX(ctx context.Context, config Config, url string, headers ...Header) (resp []byte, err error) {
	client := getClient(config)
	url = strings.TrimSpace(url)
	begin := time.Now()
	defer func() {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"url":  url,
			"time": time.Since(begin).Milliseconds(),
			"resp": summary(resp),
			"err":  err,
		}).Trace()
	}()

	return tryAgain(TryTimes, func() (resp []byte, err error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range headers {
			req.Header.Add(v.Key, v.Value)
		}
		httpResp, err := client.Do(req)
		defer func() {
			if httpResp != nil && httpResp.Body != nil {
				_ = httpResp.Body.Close()
			}
		}()
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(httpResp.Body)
	})
}

func PostX(ctx context.Context, config Config, url string, data []byte, contentType string, headers ...Header) (resp []byte, err error) {
	return commonPost(ctx, getClient(config), url, data, contentType, headers...)
}

func getClient(config Config) *http.Client {
	locker.RLock()
	client, ok := pool[config]
	locker.RUnlock()
	if !ok {
		locker.Lock()
		if client, ok = pool[config]; !ok {
			client = &http.Client{Transport: Transport(config)}
			pool[config] = client
		}
		locker.Unlock()
	}
	return client
}

func commonPost(ctx context.Context, client *http.Client, url string, data []byte, contentType string, headers ...Header) (resp []byte, err error) {
	begin := time.Now()
	defer func() {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"url":          url,
			"time":         time.Since(begin).Milliseconds(),
			"context-type": contentType,
			"post-data":    summary(data),
			"resp":         summary(resp),
			"err":          err,
		}).Trace()
	}()

	return tryAgain(TryTimes, func() (resp []byte, err error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
		for _, v := range headers {
			req.Header.Add(v.Key, v.Value)
		}
		httpResp, err := client.Do(req)
		defer func() {
			if httpResp != nil && httpResp.Body != nil {
				_ = httpResp.Body.Close()
			}
		}()
		if err != nil {
			return nil, err
		}

		return ioutil.ReadAll(httpResp.Body)
	})
}

func summary(data []byte) string {
	if len(data) > 10240 {
		return string(data[:10240]) + "..."
	}
	return *(*string)(unsafe.Pointer(&data))
}
