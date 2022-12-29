package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaConf struct {
	BrokerList string
	Topic      string
}

type ErrorLog struct {
	Kafka KafkaConf

	Url             string `json:"url"`
	Referer         string `json:"referer"`
	AccessIp        string `json:"access_ip"`
	LogType         string `json:"log_type"`
	LogTime         string `json:"log_time"`
	Level           string `json:"level"` //error、warnning
	Category        string `json:"category"`
	Text            string `json:"text"`
	ServerIp        string `json:"server_ip"`
	Host            string `json:"host"`
	AcccessQueryGet string `json:"acccess_query_get"`
}

const (
	BaseFormat       string = "2006-01-02 15:04:05"
	ERR_LOG_CHAN_LEN int    = 100000
)

var (
	ErrLogListChan chan string
	Producter      sarama.AsyncProducer
	err            error
)

// 因init函数先于main函数执行，初始化又依赖调用方的数据输入，所以不能直接用init
func initDelay(kafkaConf KafkaConf) {
	// 日志体存放chan
	ErrLogListChan = make(chan string, ERR_LOG_CHAN_LEN)

	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V1_1_1_0
	//60秒刷新一次metadata
	config.Metadata.RefreshFrequency = 30 * time.Second
	Producter, err = sarama.NewAsyncProducer(strings.Split(kafkaConf.BrokerList, ","), config)
	if err != nil {
		log.Println(err)
	}

	go sendErrLogToKafkaProccess(kafkaConf.Topic)
}

func (errLog ErrorLog) WriteErrToKafka() (err error) {
	initDelay(errLog.Kafka)
	errLog.LogTime = time.Now().Format(BaseFormat)
	hostName, _ := os.Hostname()
	errLog.ServerIp = hostName

	logBody, err := json.Marshal(errLog)
	if err != nil {
		errInfo := fmt.Sprintf("Log json formatting error: %s \n", err)
		return errors.New(errInfo)
	}

	//日志写入kafka channel, LogListChan 达到最大值的时候丢弃
	if len(ErrLogListChan) < ERR_LOG_CHAN_LEN {
		ErrLogListChan <- string(logBody)
	} else {
		return errors.New("The queue is full and the message is discarded")
	}

	return nil
}

// 异步发送错误日志消息
func sendErrLogToKafkaProccess(topic string) {
	var logBody string
	for {
		select {
		case logBody = <-ErrLogListChan:
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(logBody),
			}
			Producter.Input() <- msg
		case <-Producter.Successes():
			// 日志发送成功
		case kafkaErr := <-Producter.Errors():
			// 日志发送失败
			log.Println("Sending a message to kafka gives an error:", kafkaErr.Error())
		}
	}

	defer Producter.AsyncClose()
}
