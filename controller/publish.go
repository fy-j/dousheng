package controller

import (
	"dousheng/config"
	"dousheng/minIO"
	"dousheng/model"
	"dousheng/mq"
	"dousheng/redisUtils"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []model.VideoInfo `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	// 从Token中获取user_id
	//claims := jwt.ExtractClaims(c)
	//uid := int(claims[identityKey].(float64))
	token := c.PostForm("token")
	uid, err := GetUserIdFromToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fmt.Println(title)
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", uid, filename)
	fileObj, err := data.Open()
	if minIO.Upload(config.Conf.Bucket.Feed, finalName, fileObj, data.Size) {
		//model.VideoAdd()
		err = mq.PublishChannel.Publish("publishExchange", "publish", true, false,
			amqp.Publishing{
				Timestamp:    time.Now(),
				DeliveryMode: amqp.Persistent, //Msg set as persistent
				ContentType:  "text/plain",
				Body: mq.StructToBytes(mq.PublishMsg{
					UserId:   uid,
					FileName: finalName,
					Title:    title,
				}),
			})
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  finalName + "uploaded fill",
		})
	}

}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// 从Token中获取user_id
	claims := jwt.ExtractClaims(c)
	uid := int(claims[identityKey].(float64))
	fmt.Println(uid)
	fromRedis, err2 := redisUtils.GetVideoInfoListFromRedis(redisUtils.PUBLISHEDLIST, uid)
	//如果redis中有缓存，直接返回
	if fromRedis != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: *fromRedis,
		})
	} else if err2 == redisUtils.Nil || err2 == nil && fromRedis == nil { //如果redis中没有，从mongo中查到数据，存到redis中，过期时间设置为15min
		id, err := model.VideoListByUserID(uid, time.Now().Unix(), 30)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		}
		redisUtils.Set(redisUtils.Generate(redisUtils.PUBLISHEDLIST, strconv.FormatInt(int64(uid), 10)),
			id, time.Minute*15)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		}
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: id,
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err2.Error()})
	}
}
