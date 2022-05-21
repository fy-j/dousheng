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
	video := DemoVideos[0]
	s := video.Encoder()
	client.Do("zadd", "feedVideos", 1, s)
	byteGet := client.ZRange("feedVideos", 0, 0)
	videoPull := Decoder(string(byteGet.Val()[0]))
	fmt.Println(video)
	fmt.Println(videoPull)
}
