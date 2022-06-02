package main

import (
	"dousheng/config"
	"dousheng/controller"
	"dousheng/minIO"
	_ "dousheng/model"
	"dousheng/mq"
	"dousheng/redisUtils"
	"dousheng/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	//读配置文件
	// 应该init自动加载
	//redis初始化
	if err := redisUtils.Init(&config.Conf.Redis); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redisUtils.Close()

	//minIO初始化
	minIO.InitMinIO(&config.Conf.MinIO)

	//ampq 初始化
	mq.InitAmqp(&config.Conf.RabbitMQ)
	defer mq.PublishChannel.Close()
	defer mq.AmqpClient.Close()
	//开启消费者
	go func() { mq.Consume() }()

	//redis feed流缓存预热
	controller.RedisDataPreLoad()

	//路由初始化
	r := gin.Default()

	router.InitRouter(r)

	r.Run()

}
