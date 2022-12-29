package util

import (
	"fmt"
	"gin_frame/pkg/setting"
	logger2 "gin_frame/pkg/util/logger"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type MyLogger struct {
	ZapLogger *zap.Logger
}

var Log *MyLogger

var LoginUidKey = "loginUid"
var LoginUserNameKey = "loginUserName"
var LoginUserRoleKey = "loginUserRole"
var LoginUserStatusKey = "loginUserStatus"
var LoginSrc = "loginSrc"

func init() {
	zapLoggerT := logger2.Logger{
		ToConsole:   setting.AppC.GetBool("log.to_console"),    // 标准输出开关
		ToKafka:     setting.AppC.GetBool("log.to_kakfa"),      // 发kafka开关
		ToLocalFile: setting.AppC.GetBool("log.to_local_file"), // 记录到文件开关
		FileName:    setting.AppC.GetString("log.file_name"),
		KafkaConf: logger2.KafkaConf{
			Producer:   nil,
			BrokerList: viper.GetString("kafka.BROKERS"),
			Topic:      viper.GetString("kafka.LOG_TOPIC"),
		},
	}
	// 生产环境，只用kafka // TODO
	if !gin.IsDebugging() {
		//zapLoggerT.ToConsole = false
		//zapLoggerT.ToLocalFile = false
		//zapLoggerT.ToKafka = true
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

	// 获取函数名、文件位置、行号等
	callerLog := GetCallerInfoForLog(2)
	caller := fmt.Sprintf("%v:%v %v", callerLog["file"], callerLog["line"], callerLog["func"])

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

		//函数名，文件位置，行号
		zap.String("caller_file_line_func", caller),
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

func (log *MyLogger) AccessInfo(c *gin.Context, msg string, latencyTime time.Duration) {
	hostName, _ := os.Hostname()
	serverIp, _ := GetServerIp()
	// 状态码
	statusCode := c.Writer.Status()

	log.ZapLogger.Info(
		msg,
		zap.String("host", c.Request.Host),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("status", strconv.Itoa(statusCode)),
		zap.Int64("duration", latencyTime.Microseconds()),
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

func Error(category, msg string) {
	hostName, _ := os.Hostname()
	serverIp, _ := GetServerIp()

	Log.ZapLogger.Error(
		msg,
		zap.String("host", ""),
		//zap.String("env", Log.Env),
		zap.String("category", category),
		zap.String("server_ip", serverIp),
		zap.String("server_hostname", hostName),
		zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
	)
}

func Info(category, msg string) {
	hostName, _ := os.Hostname()
	serverIp, _ := GetServerIp()

	Log.ZapLogger.Info(
		msg,
		zap.String("host", ""),
		//zap.String("env", Log.Env),
		zap.String("category", category),
		zap.String("server_ip", serverIp),
		zap.String("server_hostname", hostName),
		zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
	)
}
