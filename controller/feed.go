package controller

import (
	"bytes"
	"dousheng/model"
	"dousheng/redisUtils"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
	client := redisUtils.Clients
	key := redisUtils.Generate("feedVideos")
	var VideoListRes []Video
	latest := c.Query("latest_time")
	latestTime := time.Now().Unix()
	if latest != "" {
		latestTime, _ = strconv.ParseInt(latest, 10, 64)
		latestTime /= 1000
	}
	//redis拉取
	if client.ZCard(key).Val() != 0 {
		//有
		//获取到序列化的字符串数组
		var tmp [30]Video
		opt := redis.ZRangeBy{
			Min: "(" + strconv.Itoa(0),
			Max: "(" + strconv.Itoa(int(latestTime)),
		}
		vs, _ := client.ZRevRangeByScore(key, opt).Result()
		videosNum := len(vs)

		//反序列化
		for pos, s := range vs {
			if pos == 30 {
				break
			}
			video := Decoder(s)
			tmp[pos] = video
		}
		if videosNum > 30 {
			VideoListRes = tmp[0:30]
		} else {
			VideoListRes = tmp[0:videosNum]
		}
	} else {
		//没有，从数据库拉取
		if InfoList, err := model.VideoList(0, 1000); err != nil {
			fmt.Println(err)
		} else {
			var tmp [1000]Video
			for pos, info := range InfoList {
				//if pos == 30 {
				//	break
				//}
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
				client.Expire(key, time.Minute*2)
			}
			if len(InfoList) > 30 {
				VideoListRes = tmp[0:30]
			} else {
				VideoListRes = tmp[0:len(InfoList)]
			}

		}
	}
	token := c.Query("token")
	if token != "" {
		//解析token得到当前用户信息
		uid, _ := GetUserIdFromToken(token)
		for pos, video := range VideoListRes {
			key = redisUtils.Generate(redisUtils.ISFACRES, strconv.FormatInt(video.Author.Id, 10), strconv.Itoa(uid))
			//查redis
			isFavRes := client.Get(key).Val()
			if isFavRes != "" {
				//有，更新video信息
				if isFavRes == "true" {
					VideoListRes[pos].IsFavorite = true
				} else if isFavRes == "false" {
					VideoListRes[pos].IsFavorite = false
				}
			} else {
				//没有，数据库查询
				res, _ := model.VideoIsFav(uid, int(video.Id))
				VideoListRes[pos].IsFavorite = res
				//更新到redis，设置ttl
				if res {
					client.Set(key, string("true"), time.Minute)
				} else {
					client.Set(key, string("false"), time.Minute)
				}
			}
			VideoListRes[pos].FavoriteCount = int64(GetFavCount(int(video.Id)))
			//关注信息
			key = redisUtils.Generate(redisUtils.ISFOLLOWED, strconv.FormatInt(video.Author.Id, 10), strconv.Itoa(uid))
			isFollowed := client.Get(key).Val()
			if isFollowed != "" {
				//有，更新video信息
				if isFollowed == "true" {
					VideoListRes[pos].Author.IsFollow = true
				} else if isFollowed == "false" {
					VideoListRes[pos].Author.IsFollow = false
				}
			} else {
				//没有，数据库查询
				res, _ := model.UserIsFollowers(uid, int(video.Author.Id))
				VideoListRes[pos].Author.IsFollow = res
				//更新到redis，设置ttl
				if res {
					client.Set(key, string("true"), time.Minute)
				} else {
					client.Set(key, string("false"), time.Minute)
				}
			}
			if uid == int(video.Author.Id) {
				VideoListRes[pos].Author.IsFollow = true
			}
		}

	}
	returnTime := int64(0)
	if len(VideoListRes) != 0 {
		videoId := VideoListRes[len(VideoListRes)-1].Id
		firstVideo, _ := model.VideoMegByID(int(videoId))
		returnTime = firstVideo.Time
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: VideoListRes,
		NextTime:  returnTime * 1000,
	})
}

func GetFavCount(videoId int) int {
	key := redisUtils.Generate(redisUtils.FAVCOUNT, strconv.Itoa(videoId))
	client := redisUtils.Clients
	count := client.Get(key).Val()
	if count == "" {
		//
		video, _ := model.VideoMegByID(int(videoId))
		count = strconv.Itoa(video.FavCount)
		client.Set(key, count, time.Minute)
	}
	res, _ := strconv.Atoi(count)
	return res
}

func RedisDataPreLoad() {
	fmt.Println("缓存预热中...")
	key := redisUtils.Generate("feedVideos")
	client := redisUtils.Clients
	if InfoList, err := model.VideoList(0, 1000); err != nil {
		fmt.Println(err)
	} else {
		var tmp [1000]Video
		for pos, info := range InfoList {
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
			client.Expire(key, time.Minute*10)
		}
	}
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
