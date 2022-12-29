package model

import (
	"gin_frame/datasource"
	"gin_frame/pkg/setting"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DbHandle    map[string]*gorm.DB
	RedisHandle map[string]*redis.Client
)

// 初始化models层需要的资源
func init() {
	initDbs()
	initRedis()
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
