package datasource

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func InstanceRedis(conf map[string]interface{}, logMode bool) *redis.Pool {
	//fmt.Println(time.Duration(conf["readTimeout"].(int64)))
	return &redis.Pool{
		MaxIdle:     int(conf["maxIdle"].(int64)),               // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
		MaxActive:   int(conf["maxActive"].(int64)),             // 最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）。
		IdleTimeout: time.Duration(conf["idleTimeout"].(int64)), // 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用。
		Wait:        true,                                       // 如果超过最大连接，是报错，还是等待。
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%v:%v", conf["host"], conf["port"]),
				redis.DialPassword(conf["pwd"].(string)),
				redis.DialConnectTimeout(time.Duration(conf["connectTimeout"].(int64))*time.Second),
				redis.DialReadTimeout(time.Duration(conf["readTimeout"].(int64))*time.Second),
				redis.DialWriteTimeout(time.Duration(conf["writeTimeout"].(int64))*time.Second),
			)
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < 2*time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
