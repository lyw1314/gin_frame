package logger

import (
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaConf struct {
	Producer sarama.AsyncProducer

	BrokerList string
	Topic      string

	//LogListChan    chan string //异步发送时的消息队列
	//LogListChanLen int         //异步发送时，消息队列长度，队列满就丢弃消息
}

func (kafkaConf *KafkaConf) NewAsyncProducer() error {
	config := sarama.NewConfig()
	//是否等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForLocal
	//config.Producer.Flush.Messages = 10 // todo 待测试
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

	producer, err := sarama.NewAsyncProducer(strings.Split(kafkaConf.BrokerList, ","), config)
	kafkaConf.Producer = producer

	//if kafkaConf.LogListChanLen <= 0 {
	//	kafkaConf.LogListChanLen = 10000
	//}
	//kafkaConf.LogListChan = make(chan string, kafkaConf.LogListChanLen)
	//go kafkaConf.SendLogToKafkaProccess()
	return err
}

func (kafkaConf *KafkaConf) Write(p []byte) (n int, err error) {
	//日志写入kafka channel, LogListChan 达到最大值的时候丢弃
	//if len(kafkaConf.LogListChan) < kafkaConf.LogListChanLen {
	//	kafkaConf.LogListChan <- string(p)
	//	return len(p), nil
	//} else {
	//	return 0, errors.New("队列已满，日志被丢弃")
	//}

	msg := &sarama.ProducerMessage{}
	msg.Topic = kafkaConf.Topic
	msg.Value = sarama.ByteEncoder(string(p))
	kafkaConf.Producer.Input() <- msg
	return len(p), nil
	//if err != nil {
	//	return
	//}
	//return
}

// 异步发送错误日志消息
//func (kafkaConf *KafkaConf) SendLogToKafkaProccess() {
//	defer kafkaConf.Producer.AsyncClose()
//
//	var logBody string
//	for {
//		select {
//		case logBody = <-kafkaConf.LogListChan:
//			msg := &sarama.ProducerMessage{
//				Topic: kafkaConf.Topic,
//				Value: sarama.ByteEncoder(logBody),
//			}
//			kafkaConf.Producer.Input() <- msg
//		case <-kafkaConf.Producer.Successes():
//			// 日志发送成功
//			fmt.Println("日志发送成功")
//		case kafkaErr := <-kafkaConf.Producer.Errors():
//			log.Println("日志发往kafka异常:", kafkaErr.Error())
//		}
//	}
//}
