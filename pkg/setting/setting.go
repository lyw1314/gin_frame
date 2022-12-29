package setting

import (
	"flag"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	RunMode        string
	GormLogMode    bool
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	HttpServerPort int

	OauthURI string

	AppC *viper.Viper
)

// 初始化配置到内存
func init() {
	var configDir = flag.String("config_dir", "./conf", "The Config Dir")
	confFile := new(string)
	var env = flag.String("env", "dev", "运行环境：dev、pro")
	flag.Parse()

	if *env == "pro" {
		*confFile = *configDir + "/prod"
		gin.SetMode(gin.ReleaseMode)
	} else {
		*confFile = *configDir + "/dev"
		gin.SetMode(gin.DebugMode)
	}

	//var confDir = flag.String("conf_dir", "./conf", "The Config Dir")
	//confFile := new(string)
	//var env = flag.String("env", "dev", "运行环境：dev、pro")
	//flag.Parse()
	//
	//if *env == "pro" {
	//	*confFile = *confDir + "/prod"
	//} else {
	//	*confFile = *confDir + "/dev"
	//}

	// 初始化配置文件
	AppC = viper.New()
	AppC.SetConfigName("app")
	AppC.AddConfigPath(*confFile)
	AppC.SetConfigType("toml")

	err := AppC.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
}

func LoadBase() {
	RunMode = AppC.GetString("RUN_MODE")
	GormLogMode = AppC.GetBool("GORM_LOG_MODE")
}

func LoadServer() {
	ReadTimeout = AppC.GetDuration("server.READ_TIMEOUT") * time.Second
	WriteTimeout = AppC.GetDuration("server.WRITE_TIMEOUT") * time.Second
	HttpServerPort = AppC.GetInt("server.HTTP_SERVER_PORT")
}

func LoadApp() {
	OauthURI = viper.GetString("oauth.HOST")
}
