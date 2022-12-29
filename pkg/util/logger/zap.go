package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	ToConsole   bool
	ToKafka     bool
	ToLocalFile bool
	KafkaConf   KafkaConf
	FileName    string
}

func (logger *Logger) NewZapLogger() (*zap.Logger, error) {
	var err error
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel
	})

	coreArr := make([]zapcore.Core, 0)
	if logger.ToConsole == true {
		consoleDebugging := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		coreArr = append(coreArr,
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))
	}

	if logger.ToKafka == true {
		err = logger.KafkaConf.NewAsyncProducer()
		topicErrors := zapcore.AddSync(&logger.KafkaConf)
		kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		coreArr = append(coreArr, zapcore.NewCore(kafkaEncoder, topicErrors, highPriority))
	}

	if logger.ToLocalFile == true {
		encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		writeSyncer := logger.getLogWriter()
		coreArr = append(coreArr, zapcore.NewCore(encoder, writeSyncer, highPriority))
	}

	zapLogger := zap.New(zapcore.NewTee(coreArr...))
	defer zapLogger.Sync()
	return zapLogger, err
}

func (logger *Logger) getLogWriter() zapcore.WriteSyncer {
	// 切割日志文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logger.FileName,
		MaxSize:    100,   //在进行切割之前，日志文件的最大大小（以MB为单位
		MaxBackups: 50,    //保留旧文件的最大个数
		MaxAge:     90,    //保留旧文件的最大天数
		Compress:   false, //是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}
