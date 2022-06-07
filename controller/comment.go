package controller

import (
	"dousheng/model"
	"dousheng/redisUtils"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []model.AssessmentInfo `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uid := int(claims[identityKey].(float64))
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	actionType, _ := strconv.Atoi(c.Query("action_type"))
	commentText := c.Query("comment_text")
	commentId, _ := strconv.Atoi(c.Query("comment_id"))
	if actionType == 1 {
		model.AssAdd(uid, videoId, commentText)
	} else {
		model.AssDel(videoId, commentId)
	}
	res, err := model.AssListByVideoID(videoId)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1, StatusMsg: err.Error()},
			CommentList: nil,
		})
	}
	fmt.Println(redisUtils.Generate(redisUtils.ASSESSMENT, strconv.Itoa(videoId), strconv.Itoa(uid)))
	err = redisUtils.Set(redisUtils.Generate(redisUtils.ASSESSMENT, strconv.Itoa(videoId), strconv.Itoa(uid)), res, time.Minute*15)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1, StatusMsg: err.Error()},
			CommentList: nil,
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response{0, ""},
		res,
	})
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uid := int(claims[identityKey].(float64))
	videoId, _ := strconv.Atoi(c.Query("video_id"))
	result, err := redisUtils.GetAssessmentFromRedis(videoId, uid)
	if err != redisUtils.Nil && err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 1, StatusMsg: err.Error()},
			CommentList: nil,
		})
	} else if err == redisUtils.Nil && result == nil {
		res, err := model.AssListByVideoID(videoId)
		if err != nil {
			c.JSON(http.StatusOK, CommentListResponse{
				Response{1, err.Error()},
				nil,
			})
		} else {
			redisUtils.Set(redisUtils.Generate(redisUtils.ASSESSMENT, strconv.Itoa(videoId), strconv.Itoa(uid)), res, time.Minute*15)
			c.JSON(http.StatusOK, CommentListResponse{
				Response{0, ""},
				res,
			})
		}
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			Response{0, ""},
			result,
		})
	}
}
