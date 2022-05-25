package router

import (
	"dousheng/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	//redis test api
	apiRouter.GET("/redis_test", controller.RedisTest)

	// basic apis
	apiRouter.GET("/feed/", controller.AuthMiddleware.MiddlewareFunc(), controller.Feed)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.AuthMiddleware.LoginHandler)
	apiRouter.GET("/user", controller.AuthMiddleware.MiddlewareFunc(), controller.UserInfoHandler)
	apiRouter.POST("/publish/action/", controller.AuthMiddleware.MiddlewareFunc(), controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)

}
