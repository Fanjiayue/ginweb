package common

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/gomodule/redigo/redis"

)

var RedisClient *redis.Pool

func InitRedis() *redis.Pool{

	RedisClient = &redis.Pool{
		MaxIdle:     viper.GetInt("redis.MaxIdle"), /*最大的空闲连接数*/
		MaxActive:   viper.GetInt("redis.MaxActive"), /*最大的激活连接数*/
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", viper.GetString("redis.host")+":"+viper.GetString("redis.port"), redis.DialPassword(""))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	fmt.Println("redis connect success")
	return RedisClient

}

func GetRedis() redis.Conn{
	return RedisClient.Get()
}


