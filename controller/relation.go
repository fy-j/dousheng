package controller

import (
	"dousheng/model"
	"dousheng/redisUtils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type UserListResponse struct {
	Response
	UserList []model.UserInfo `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	//userId, _ := strconv.Atoi(c.Query("user_id"))
	toUserId, _ := strconv.Atoi(c.Query("to_user_id"))
	uid, err := GetUserIdFromToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
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
	//		StatusMsg:  "token验证失败：",
	//	})
	//	println(uid, userId)
	//	return
	//}
	err = model.UserFollow(uid, toUserId, actionType)
	//同步redis
	key := redisUtils.Generate(redisUtils.ISFACRES, strconv.Itoa(toUserId), strconv.Itoa(uid))
	if actionType == 1 {
		redisUtils.Clients.Set(key, string("true"), time.Minute)
	} else {
		redisUtils.Clients.Set(key, string("false"), time.Minute)
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
		StatusMsg:  "成功follow",
	})
	//token := c.Query("token")
	//
	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 0})
	//} else {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userId, _ := strconv.Atoi(c.Query("user_id"))
	uid, err := GetUserIdFromToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	if uid != userId {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "token验证失败",
		})
		return
	}
	list, err := model.UserFollowList(uid)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: nil,
		})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功获取关乎列表",
		},
		UserList: list,
	})
	//c.JSON(http.StatusOK, UserListResponse{
	//	Response: Response{
	//		StatusCode: 0,
	//	},
	//	UserList: []User{DemoUser},
	//})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	uid, err := GetUserIdFromToken(token)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: nil,
		})
	}
	//userId, _ := strconv.Atoi(c.Query("user_id"))
	//println(userId, uid)
	//if uid != userId {
	//	c.JSON(http.StatusOK, UserListResponse{
	//		Response: Response{
	//			StatusCode: 1,
	//			StatusMsg:  "token验证失败",
	//		},
	//		UserList: nil,
	//	})
	//}
	userList, err := model.UserFansList(uid)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "成功得到粉丝信息",
		},
		UserList: userList,
	})
}
