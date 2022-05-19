package main

import (
	"dousheng/config"
	"dousheng/minIO"
	"dousheng/redis"
	"dousheng/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	//读配置文件
	config.Init()

	//redis初始化
	if err := redis.Init(&config.Conf.Redis); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	//minIO初始化
	minIO.InitMinIO(&config.Conf.MinIO)

	//路由初始化
	r := gin.Default()

	router.InitRouter(r)

	r.Run()

}
