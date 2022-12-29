package main

import (
	"fmt"
	logger2 "gin_frame/pkg/util/logger"
	"time"

	"go.uber.org/zap"
)

func main() {
	var zapLoggerT = logger2.Logger{
		ToConsole: true,
		ToKafka:   true,
		KafkaConf: logger2.KafkaConf{
			Producer:   nil,
			BrokerList: "test-xxx.net:9092",
			Topic:      "dj.server.log",
			//LogListChanLen: 2,
		},
	}
	zapLogger, err := zapLoggerT.NewZapLogger()
	if err != nil {
		fmt.Println(err.Error())
	}

	for i := 0; i < 5; i++ {
		go func() {
			zapLogger.Info("test log info",
				zap.String("host", "test"),
				zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
			)
		}()
	}

	time.Sleep(2 * time.Second)
	zapLogger.Info("test11 log info",
		zap.String("host", "test"),
		zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
	)

	//zapLogger.Error("test log error",
	//	zap.String("host", "xxx.cn"),
	//	zap.String("log_time", time.Now().Format("2006-01-02 15:04:05")),
	//)
	select {}
}
