package controller

import (
	"dousheng/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	userId, _ := strconv.Atoi(c.Query("user_id"))
	videoId, err := strconv.Atoi(c.Query("video_id"))
	uid, err := GetUserIdFromToken(token)
	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	if uid != userId {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token伪造",
		})
	}
	model.VideoFavAction(uid, videoId, actionType)
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		//VideoList: DemoVideos,
	})
}
