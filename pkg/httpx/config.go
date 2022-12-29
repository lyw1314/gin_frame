// author by lipengfei5 @2022-03-24

package httpx

var defaultConfig = Config{
	DialTimeout:         3000,
	DialKeepAlive:       30000,
	MaxIdleConnection:   100,
	MaxIdleConnsPerHost: 100,
	IdleConnTimeout:     30000,
	EnableHttp2:         true,
	InsecureSkipVerify:  true,
}
var TryTimes = 10

type Config struct {
	DialTimeout         int64 // ms
	DialKeepAlive       int64 // ms
	MaxIdleConnection   int
	MaxIdleConnsPerHost int
	IdleConnTimeout     int64 // ms
	EnableHttp2         bool
	InsecureSkipVerify  bool
}

type Header struct {
	Key   string
	Value string
}
