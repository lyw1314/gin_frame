package datasource

import (
	"fmt"
	"gin_frame/pkg/util"
	"strconv"
	"time"

	"github.com/beltran/gohive"
	_ "github.com/go-sql-driver/mysql"
)

var (
	HiveConn *gohive.Connection
)

type HiveConfig struct {
	FetchSize int64
	Username  string
	Password  string
	Host      string
	Port      int
	Connect   *gohive.Connection
}

func reportToMonitorSystem(brokers, detail string, err error) {
	category := "sdk2 ba report sync failed"
	//往kafka写错误日志
	errLog := util.ErrorLog{
		Kafka: util.KafkaConf{BrokerList: brokers,
			Topic: "djplat.shbt.error"},
		Host:            "gocron",
		Url:             "cron",
		Category:        category,
		Text:            category + " " + detail + "\n" + err.Error(),
		LogType:         "customError",
		Referer:         "cron",
		AccessIp:        "cron",
		Level:           "error",
		AcccessQueryGet: "cron",
	}

	if res := errLog.WriteErrToKafka(); res != nil {
		util.Error("reportToMonitorSystem", res.Error())
	} else {
		util.Info("reportToMonitorSystem", "Error has been sent to kafka")
	}
}

func InstanceHive(brokers string, conf map[string]interface{}, logMode bool) *gohive.Connection {

	portStr := fmt.Sprintf("%v", conf["port"])
	port, _ := strconv.Atoi(portStr)

	configuration := gohive.NewConnectConfiguration()
	configuration.Service = "hive"
	configuration.FetchSize = 10000
	configuration.Username = fmt.Sprintf("%v", conf["user"])
	configuration.Password = fmt.Sprintf("%v", conf["pwd"])
	configuration.TransportMode = "binary"

	hiveConn, errConn := gohive.Connect(fmt.Sprintf("%v", conf["host"]), port, "LDAP", configuration)
	if errConn != nil {
		reportToMonitorSystem(brokers, "InstanceHive err", errConn)
		//当链接遇到错误的时候，进行重试10次
		for i := 1; i <= 10; i++ {
			reportToMonitorSystem(brokers, "connect hive error. reconnect... times:"+strconv.Itoa(i), errConn)
			//log.Printf("reconnect... times %d",i)
			hiveConn, errConn = gohive.Connect(fmt.Sprintf("%v", conf["host"]), port, "LDAP", configuration)
			if errConn == nil {
				break
			}
			// 1000ms后重试第二次
			time.Sleep(1000 * time.Millisecond)
		}
		return nil
	}
	//defer HiveConn.Close()
	return hiveConn
}
