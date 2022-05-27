package mq

import (
	"dousheng/minIO"
	"dousheng/model"
	"dousheng/redisUtils"
	"fmt"
	"log"
	"strconv"
	"time"
)

func Consume() {
	deliveries, err := PublishChannel.Consume("publishQueue", "any", false, false, false, true, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	v, ok := <-deliveries
	if ok {
		data := BytesToStruct(v.Body).(PublishMsg)
		fmt.Println("收到消息", data)
		url := minIO.GetURL(data.FileName, time.Second*24*60*60)
		fmt.Println(url)
		model.VideoAdd(data.UserId, "", url, data.Title)
		//redis删除
		client := redisUtils.Clients
		client.Del(redisUtils.Generate(redisUtils.PUBLISHEDLIST, strconv.FormatInt(int64(data.UserId), 10)))
		client.Del(redisUtils.Generate("feedVideos"))
		if err := v.Ack(true); err != nil {
			fmt.Println(err.Error())
		}
	}
}
