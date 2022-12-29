package datasource

import (
	"time"

	"github.com/beltran/gohive"
)

var (
	HiveConnect *gohive.Connection
)

type HiveConfig2 struct {
	FetchSize int64
	Username  string
	Password  string
	Host      string
	Port      int
	Connect   *gohive.Connection
}

func NewHive(c HiveConfig2) (*HiveConfig2, error) {
	if c.FetchSize == 0 {
		c.FetchSize = 1000
	}
	if c.Port == 0 {
		c.Port = 10000
	}

	h := &HiveConfig2{
		FetchSize: c.FetchSize,
		Username:  c.Username,
		Password:  c.Password,
		Host:      c.Host,
		Port:      c.Port,
	}
	hiveConn, errConn := h.connect()
	if errConn != nil {
		return nil, errConn
	}

	h.Connect = hiveConn
	HiveConnect = hiveConn

	return h, nil
}

func (c *HiveConfig2) connect() (*gohive.Connection, error) {
	//start := time.Now()
	configuration := gohive.NewConnectConfiguration()
	configuration.Service = "hive"
	configuration.FetchSize = c.FetchSize
	configuration.Username = c.Username
	configuration.Password = c.Password
	configuration.TransportMode = "binary"

	hiveConn, errConn := gohive.Connect(c.Host, c.Port, "LDAP", configuration)
	if errConn != nil {
		// 当链接遇到错误的时候，进行重试10次
		for i := 1; i <= 10; i++ {
			hiveConn, errConn = gohive.Connect(c.Host, c.Port, "LDAP", configuration)
			if errConn == nil {
				return hiveConn, nil
			}

			// 1000ms后重试第二次
			time.Sleep(1000 * time.Millisecond)
		}
		return nil, errConn
	}

	// WriteLog(datamodels.Log{
	// 	LogLevel: INFO,
	// 	LogTag:   "CONNECT_HIVE",
	// 	LogInfo:  "connect hive success",
	// }, "costTime:", time.Since(start))

	//defer HiveConn.Close()
	return hiveConn, nil
}

func (c *HiveConfig2) Conn() *gohive.Connection {
	return c.Connect
}

func (c *HiveConfig2) Clone() *gohive.Cursor {
	return c.Connect.Cursor()
}

func (c *HiveConfig2) Close() {
	c.Connect.Close()
}

func GetHiveConfig(conf map[string]interface{}) HiveConfig2 {
	//conf = model.Config["hiveconfig.hdp_ads_sdkrec"].(map[string]interface{})
	return HiveConfig2{
		FetchSize: 0,
		Username:  conf["user"].(string),
		Password:  conf["pwd"].(string),
		Host:      conf["host"].(string),
		Port:      int(conf["port"].(int64)),
	}
}

func HiveSetup(conf map[string]interface{}) (*HiveConfig2, error) {
	return NewHive(GetHiveConfig(conf))
}
