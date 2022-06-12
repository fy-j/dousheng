package controller

import (
	"dousheng/model"
	"dousheng/redisUtils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	//userId, _ := strconv.Atoi(c.Query("user_id"))
	videoId, err := strconv.Atoi(c.Query("video_id"))
	uid, err := GetUserIdFromToken(token)
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//if uid != userId {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  "token伪造",
	//	})
	//}
	err = model.VideoFavAction(uid, videoId, actionType)
	//同步到redis
	key1 := redisUtils.Generate(redisUtils.ISFACRES, strconv.Itoa(videoId), strconv.Itoa(uid))
	key2 := redisUtils.Generate(redisUtils.FAVCOUNT, strconv.Itoa(videoId))
	if actionType == 1 {
		redisUtils.Clients.Set(key1, string("true"), time.Minute)
		redisUtils.Clients.Del(key2)
	} else {
		redisUtils.Clients.Set(key1, string("false"), time.Minute)
		redisUtils.Clients.Del(key2)
	}
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "成功喜欢",
	})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uid := int(claims[identityKey].(float64))
	userId, err1 := strconv.Atoi(c.Query("user_id"))
	if err1 != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err1.Error(),
			},
			VideoList: nil,
		})
	}
	fromRedis, err := redisUtils.GetVideoInfoListFromRedis(redisUtils.ISFAVORITE, uid)
	if fromRedis != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: *fromRedis,
		})
	} else if err == redisUtils.Nil || err == nil && fromRedis == nil {
		videoList, err := model.VideoFavList(userId)
		if err != nil {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err1.Error(),
				},
				VideoList: nil,
			})
		}
		redisUtils.Set(redisUtils.Generate(redisUtils.ISFAVORITE, strconv.Itoa(userId)), videoList, time.Minute*15)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: videoList,
		})
	}
}
