package logger

type logStruct struct {
	RequestID      string `json:"request_id"` //请求唯一标识
	Host           string `json:"host"`       //域名，必填
	Scheme         string `json:"scheme"`     //协议
	Method         string `json:"method"`     //http请求方法
	Path           string `json:"path"`       //http请求路径
	StatusCode     string `json:"status"`     //http请求状态码
	Duration       string `json:"duration"`   //请求处理持续时间
	Level          string `json:"level"`      //日志等级：debug、info、warning、error、panic、fatal
	Category       string `json:"category"`   //log摘要
	Msg            string `json:"msg"`        //log详情
	UserID         string `json:"uid"`
	AccessIP       string `json:"access_ip"`
	ServerIP       string `json:"server_ip"`
	ServerHostname string `json:"server_hostname"`
	LogTime        string `json:"log_time"`
	UserAgent      string `json:"user_agent"`
}
