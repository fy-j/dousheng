package controller

import (
	"bytes"
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
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: DemoVideos,
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
func Decoder(videoByte string) Video {
	var video Video
	decoder := gob.NewDecoder(bytes.NewReader([]byte(videoByte)))
	err := decoder.Decode(&video)
	if err != nil {
		log.Fatal(err)
	}
	return video
}
