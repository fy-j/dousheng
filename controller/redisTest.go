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
	//for _, video := range DemoVideos {
	//	s := video.Encoder()
	//	client.Do("zadd", "feedVideos", time.Now().Unix(), s)
	//}
	//测试数据
	//model.VideoAdd(1, "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", "https://www.w3schools.com/html/movie.mp4", "aaa")
	//time.Sleep(time.Second)
	//model.VideoAdd(1, "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", "https://www.w3schools.com/html/movie.mp4", "bbb")
	//time.Sleep(time.Second)
	//model.VideoAdd(1, "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", "https://www.w3schools.com/html/movie.mp4", "ccc")
	//time.Sleep(time.Second)
	//model.VideoAdd(1, "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", "https://www.w3schools.com/html/movie.mp4", "ddd")
	//time.Sleep(time.Second)
	//model.VideoAdd(1, "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", "https://www.w3schools.com/html/movie.mp4", "eee")
	//time.Sleep(time.Second)

}
