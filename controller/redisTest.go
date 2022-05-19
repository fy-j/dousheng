package controller

import (
	"dousheng/redis"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RedisTest(c *gin.Context) {
	client := redis.Clients
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	//client.FlushAll()

	val, err := client.Get("key").Result()
	if err == redis.Nil {
		fmt.Println("key does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key", val)
	}
}
