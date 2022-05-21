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
	//client.FlushAll()
	for _, video := range DemoVideos {
		s := video.Encoder()
		client.Do("zadd", "feedVideos", 1, s)
	}
}
