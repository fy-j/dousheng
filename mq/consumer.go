package mq

import (
	"bytes"
	"dousheng/config"
	"dousheng/minIO"
	"dousheng/model"
	"dousheng/redisUtils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
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
		coverPath := CreateCoverPath(data.FileName)
		CreateCover(url, coverPath)
		buf := new(bytes.Buffer)
		file, _ := os.Open(config.Conf.Video.Address + coverPath + "-cover.jpg")
		buf.ReadFrom(file)
		if err := minIO.Upload(config.Conf.Bucket.Feed, coverPath+"-cover.jpg", buf, int64(buf.Len())); !err {
			return
		}
		coverurl := minIO.GetCoverURL(coverPath+"-cover.jpg", time.Second*120)
		fmt.Println(coverurl)
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

func CmdInit(url string, FileName string) []string {
	ls := []string{
		"-i",
		url,
		"-y",
		"-f",
		"image2",
		"-vframes",
		"1",
		config.Conf.Video.Address + FileName + "-cover" + ".jpg",
	}
	return ls
}

func CreateCoverPath(FileName string) string {
	filenameWithSuffix := path.Base(FileName)
	fmt.Println("filenameWithSuffix =", filenameWithSuffix)
	fileSuffix := path.Ext(filenameWithSuffix)
	fmt.Println("fileSuffix =", fileSuffix)
	coverPath := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return coverPath
}

func CreateCover(url string, coverPath string) {
	cmdArgs := CmdInit(url, coverPath)
	cmd := exec.Command("ffmpeg", cmdArgs...)
	if err1 := cmd.Run(); err1 != nil {
		fmt.Println(err1)
		panic("could not generate frame")
	}
}
