package mq

import (
	"fmt"
	"github.com/go-redis/redis"
)

var Rdb *redis.Client

func InitRedis(){
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
		PoolSize: 10,		//连接数
	})
	_,err := Rdb.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return
	}
}