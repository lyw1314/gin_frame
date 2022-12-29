package model

import (
	"context"
	"gin_frame/datasource"
	pb "gin_frame/grpc_pb/api_data"
	"gin_frame/pkg/setting"
	"gin_frame/pkg/util"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	DbHandle    map[string]*gorm.DB
	RedisHandle map[string]*redis.Client

	ApiDataClient pb.ApiDataClient
)

// 初始化models层需要的资源
func init() {
	initDbs()
	initRedis()

	// 初始化grpc client
	InitApiDataClient()
}

// 初始化所有DB
func initDbs() {
	DbHandle = make(map[string]*gorm.DB)
	conf := setting.AppC.GetStringMap("database")
	for k, c := range conf {
		cc := c.(map[string]interface{})
		DbHandle[k] = datasource.InstanceDb(cc, setting.GormLogMode)
	}
}

// 初始化redis
func initRedis() {
	RedisHandle = make(map[string]*redis.Client)
	conf := setting.AppC.GetStringMap("redis")
	for k, _ := range conf {
		cc := conf[k].(map[string]interface{})
		RedisHandle[k] = datasource.IniRedis(cc)
	}
}

// InitApiDataClient 初始化apiData GRPC服务
func InitApiDataClient() {
	port := setting.AppC.GetString("api_data.port")
	timeOut := time.Duration(setting.AppC.GetInt64("api_data.timeout"))
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*timeOut)
	defer cancel()
	conn, err := grpc.DialContext(ctx, port, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)))
	if err != nil {
		util.Error("category_name", err.Error())
	}

	ApiDataClient = pb.NewApiDataClient(conn)
}
