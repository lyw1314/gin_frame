package datasource

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func IniRedis(conf map[string]interface{}) *redis.Client {
	// 连接服务器
	opt := redis.Options{
		Addr:         fmt.Sprintf("%v:%v", conf["host"], conf["port"]),
		Password:     conf["pwd"].(string),
		DB:           0,
		ReadTimeout:  time.Duration(conf["read_timeout"].(int64)) * time.Second,
		WriteTimeout: time.Duration(conf["write_timeout"].(int64)) * time.Second,
		IdleTimeout:  time.Duration(conf["idle_timeout"].(int64)) * time.Second,
		DialTimeout:  time.Duration(conf["connect_timeout"].(int64)) * time.Second,
		MinIdleConns: int(conf["min_idle_cons"].(int64)),
		PoolSize:     int(conf["pool_size"].(int64)),
	}

	redisConn := redis.NewClient(&opt)
	return redisConn
}
