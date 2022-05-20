package controller

import (
	"dousheng/model"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 其他功能暂时依赖该项,待删除
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// 用户登录/注册请求
type UserRequest struct {
	Name string `json:"username"`
	Pwd  string `json:"password"`
}

// 用户登录/注册响应
type UserLoginResponse struct {
	Response
	UserId int    `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

// 用户信息响应
type UserInfoResponse struct {
	Response
	Userinfo model.UserInfo `json:"user"`
}

// Payload结构体
type Payload struct {
	UserId int
}

// MapClaims 默认key
var identityKey = "id"

// Auth中间件
var AuthMiddleware, _ = jwt.New(
	&jwt.GinJWTMiddleware{
		Realm:            "douyin",             //标识
		SigningAlgorithm: "HS256",              //加密算法
		Key:              []byte("summerCamp"), //密钥
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		IdentityKey:      identityKey,
		// 1. 用户登录流
		// 1.1 登录验证
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// 获取用户登录请求
			var loginReq UserRequest
			if err := c.ShouldBind(&loginReq); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginReq.Name
			password := loginReq.Pwd

			// 登录验证,调用Model层UserLogin函数，返回用户信息
			user_full, err := model.UserLogin(username, password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			// 通过gin.context传递最终Response中的user_id
			c.Set("uid", user_full.UserId)
			// 返回用户信息到PayloadFunc
			return user_full, nil
		},
		// 1.2 添加Payload
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// 接收Authenticator的返回值data,并将待存放的Payload字段加入MapClaims中并返回
			if v, ok := data.(model.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserId,
				}
			}
			return jwt.MapClaims{}
		},
		// 1.3 UserLoginResponse响应封装
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			uid, _ := c.Get("uid")
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0, StatusMsg: "success"},
				UserId:   uid.(int),
				Token:    message,
			})
		},
		// 2. Token鉴权流
		// 2.1 解析JWT,并将Payload传递给Authorizator
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &Payload{
				UserId: int(claims[identityKey].(float64)),
			}
		},
		// 2.2 Token鉴权,当用户通过token请求受限接口时,会执行这段逻辑
		Authorizator: func(data interface{}, c *gin.Context) bool {
			uid, err := strconv.Atoi(c.Query("user_id"))
			if err != nil {
				return false
			}
			if v, ok := data.(*Payload); ok && v.UserId == uid {
				return true
			}
			return false
		},
		// 2.3 Token鉴权失败处理
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// 指定从哪里获取token
		// 格式："<source>:<name>"
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

// Token鉴权成功,用户信息请求逻辑
func UserInfoHandler(c *gin.Context) {
	// 从Token中获取user_id
	claims := jwt.ExtractClaims(c)
	uid := int(claims[identityKey].(float64))
	// 根据user_id查询信息
	if userInfo, err := model.UserInfoById(uid); err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			Userinfo: model.UserInfo{},
		})
	} else {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 0, StatusMsg: "success"},
			Userinfo: userInfo,
		})
	}
}

// 用户注册逻辑
func Register(c *gin.Context) {
	// 接收注册信息
	var registerReq UserRequest
	if err := c.ShouldBind(&registerReq); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	name := registerReq.Name
	pwd := registerReq.Pwd

	// 有效性判断
	if utf8.RuneCountInString(name) > 32 || utf8.RuneCountInString(pwd) > 32 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "username or password exceeds 32 characters"},
		})
	} else if exist, _ := model.UserExist(name); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		// 用户信息持久化
		if _, err := model.UserAdd(name, pwd); err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "user add error",
			})
			return
		}
		// 请求重定向(307),进行用户登录
		c.Redirect(http.StatusTemporaryRedirect, "/douyin/user/login")
	}
}
