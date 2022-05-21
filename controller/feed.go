package controller

import (
	"bytes"
	"dousheng/redis"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	client := redis.Clients
	key := "feedVideos"
	var VideoList []Video
	//redis拉取
	if videosNum := client.ZCard(key).Val(); videosNum != 0 {
		//有
		//获取到序列化的字符串数组
		var tmp [30]Video
		Vs := client.ZRange(key, 0, videosNum-1).Val()
		//反序列化
		for pos, s := range Vs {
			video := Decoder(s)
			tmp[pos] = video
		}
		VideoList = tmp[0:videosNum]
	} else {
		//TODO 没有，从数据库拉取

		//TODO 更新到redis

	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: VideoList,
		NextTime:  time.Now().Unix(),
	})
}

//Video序列化
func (v *Video) Encoder() string {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(v)
	if err != nil {
		log.Fatal(err)
	}
	return string(buffer.Bytes())
}

//Video反序列化
func Decoder(videoString string) Video {
	var video Video
	decoder := gob.NewDecoder(bytes.NewReader([]byte(videoString)))
	err := decoder.Decode(&video)
	if err != nil {
		log.Fatal(err)
	}
	return video
}
