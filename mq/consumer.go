package mq

import (
	"bytes"
	"dousheng/config"
	"dousheng/minIO"
	"dousheng/model"
	"dousheng/redisUtils"
	"fmt"
	"log"
	"os/exec"
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
		cmd := exec.Command("ffmpeg", "-i", "\""+url+"\"", "-f", "image2", "-frames:v", "1", "\"D:\\"+data.FileName+"-cover\"")
		fmt.Println(url)
		buf := new(bytes.Buffer)
		cmd.Stdout = buf
		if cmd.Run() != nil {
			panic("could not generate frame")
		}
		minIO.Upload(config.Conf.Bucket.Feed, data.FileName+"—cover", buf, int64(buf.Len()))
		coverurl := minIO.GetCoverURL(data.FileName+"-cover", time.Second*120)
		model.VideoAdd(data.UserId, coverurl, url, data.Title)
		//redis删除
		client := redisUtils.Clients
		client.Del(redisUtils.Generate(redisUtils.PUBLISHEDLIST, strconv.FormatInt(int64(data.UserId), 10)))
		client.Del(redisUtils.Generate("feedVideos"))
		if err := v.Ack(true); err != nil {
			fmt.Println(err.Error())
		}
	}
}
