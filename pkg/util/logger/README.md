## 使用样例
```go 
package main

import (
	"gin_frame/pkg/logger"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func main() {
	var zapLoggerT = logger.Logger{
		ToConsole: true,    //控制日志是否打印到标准输出
		ToKafka:   true,    //控制日志是否发送到kafka
		KafkaConf: logger.KafkaConf{
			Producer:   nil,    //写日志时kafka生产者
			BrokerList: "test-kf1.adsys.shbt2.qihoo.net:9092,test-kf2.adsys.shbt2.qihoo.net:9092,test-kf3.adsys.shbt2.qihoo.net:9092",  //kafka集群broker
			Topic:      "dj.mdsp.mesos.server.log", //写日志时kafka的topic
		},
	}
	zapLogger, err := zapLoggerT.NewZapLogger()
	if err != nil {
		fmt.Println(err.Error())
	}
	zapLogger.Info("test log info",
		zap.String("host", "e.360.cn"),
		zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
		)
	zapLogger.Error("test log error",
		zap.String("host", "mobile.e.360.cn"),
		zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
	)
	select {}
}

```

## gin项目使用样例
```go
package util

import (
	"gin_frame/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

type MyLogger struct {
	ZapLogger *zap.Logger
}

var Log *MyLogger

func init() {
	zapLoggerT := logger.Logger{
		ToConsole: true,
		ToKafka:   true,
		KafkaConf: logger.KafkaConf{
			Producer:   nil,
			BrokerList: viper.GetString("kafka.BROKERS"),
			Topic:      viper.GetString("kafka.LOG_TOPIC"),
		},
	}
	if !gin.IsDebugging() {
		zapLoggerT.ToConsole = false
		zapLoggerT.ToKafka = true
	}

	zapLogger, err := zapLoggerT.NewZapLogger()
	if err != nil {
		zapLogger.Error(err.Error())
	}
	Log = &MyLogger{}
	Log.ZapLogger = zapLogger
}

func (log *MyLogger) Error(c *gin.Context, msg string) {
	hostName, _ := os.Hostname()
	serverIp, _ := GetServerIp()

	log.ZapLogger.Error(
		msg,
		zap.String("host", c.Request.Host),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("access_ip", c.ClientIP()),
		zap.String("server_ip", serverIp),
		zap.String("server_hostname", hostName),
		zap.String("uid", strconv.FormatInt(c.GetInt64(LoginUidKey), 10)),
		zap.String("request_id", c.Writer.Header().Get("X-Request-Id")),
		zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
	)
}

func (log *MyLogger) Info(c *gin.Context, msg string) {
	hostName, _ := os.Hostname()
	serverIp, _ := GetServerIp()

	log.ZapLogger.Info(
		msg,
		zap.String("host", c.Request.Host),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("access_ip", c.ClientIP()),
		zap.String("server_ip", serverIp),
		zap.String("server_hostname", hostName),
		zap.String("uid", strconv.FormatInt(c.GetInt64(LoginUidKey), 10)),
		zap.String("request_id", c.Writer.Header().Get("X-Request-Id")),
		zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
	)
}

func (log *MyLogger) Warning(c *gin.Context, msg string) {
	hostName, _ := os.Hostname()
	serverIp, _ := GetServerIp()

	log.ZapLogger.Warn(
		msg,
		zap.String("host", c.Request.Host),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("access_ip", c.ClientIP()),
		zap.String("server_ip", serverIp),
		zap.String("server_hostname", hostName),
		zap.String("uid", strconv.FormatInt(c.GetInt64(LoginUidKey), 10)),
		zap.String("request_id", c.Writer.Header().Get("X-Request-Id")),
		zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
	)
}

```
