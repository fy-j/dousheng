package controller

import (
	"bytes"
	"dousheng/model"
	"dousheng/redis"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
	key := redis.Generate("feedVideos")
	var VideoListRes []Video
	//redis拉取
	if videosNum := client.ZCard(key).Val(); videosNum != 0 {
		//有
		//获取到序列化的字符串数组
		var tmp [30]Video
		Vs := client.ZRevRange(key, 0, videosNum-1).Val()
		//反序列化
		for pos, s := range Vs {
			video := Decoder(s)
			tmp[pos] = video
		}
		VideoListRes = tmp[0:videosNum]

	} else {
		//没有，从数据库拉取
		if InfoList, err := model.VideoList(0, 30); err != nil {
			fmt.Println(err)
		} else {
			var tmp [30]Video
			for pos, info := range InfoList {
				if pos == 30 {
					break
				}
				tmp[pos] = Video{
					Id: int64(info.VideoID),
					Author: User{
						Id:            int64(info.Author.UserId),
						Name:          info.Author.Name,
						FollowCount:   int64(info.Author.FollCount),
						FollowerCount: int64(info.Author.FansCount),
						IsFollow:      info.Author.IsFollow,
					},
					PlayUrl:       info.PlayUrl,
					CoverUrl:      info.CoverUrl,
					FavoriteCount: int64(info.FavCount),
					CommentCount:  int64(info.ComCount),
					IsFavorite:    info.IsFav,
					Title:         info.Title,
				}
				//更新到redis
				client.Do("zadd", key, info.Time, tmp[pos].Encoder())
			}
			VideoListRes = tmp[0:len(InfoList)]
		}
	}
	token := c.Query("token")
	if token != "" {
		//解析token得到当前用户信息(有bug)
		uid, _ := GetUserIdFromToken(token)
		for _, video := range VideoListRes {
			key = redis.Generate(redis.ISFACRES, strconv.FormatInt(video.Id, 10), strconv.Itoa(uid))
			//查redis
			isFavRes := client.Get(key).Val()
			if isFavRes != "" {
				//有，更新video信息
				if isFavRes == "true" {
					video.IsFavorite = true
				} else if isFavRes == "false" {
					video.IsFavorite = false
				}
			} else {
				//没有，数据库查询
				res, _ := model.VideoIsFav(uid, int(video.Id))
				video.IsFavorite = res
				//更新到redis，设置ttl
				if res {
					client.Set(key, string("true"), time.Minute)
				} else {
					client.Set(key, string("false"), time.Minute)
				}
			}
			//关注信息
			key = redis.Generate(redis.ISFOLLOWED, strconv.FormatInt(video.Author.Id, 10), strconv.Itoa(uid))
			isFollowed := client.Get(key).Val()
			if isFollowed != "" {
				//有，更新video信息
				if isFollowed == "true" {
					video.Author.IsFollow = true
				} else if isFollowed == "false" {
					video.Author.IsFollow = false
				}
			} else {
				//没有，数据库查询
				res, _ := model.VideoAuthorIsFollowed(uid, int(video.Id))
				video.Author.IsFollow = res
				//更新到redis，设置ttl
				if res {
					client.Set(key, string("true"), time.Minute)
				} else {
					client.Set(key, string("false"), time.Minute)
				}
			}
		}

	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: VideoListRes,
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
