package minIO

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"time"
)

func Download(c *gin.Context) {
	reqParams := make(url.Values)
	presignedURL, err := Client.PresignedGetObject("userfeed", "shipin1", time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		c.JSON(0, gin.H{
			"msg": "获取失败",
		})
	}
	c.JSON(200, gin.H{
		"msg": "success",
		"url": fmt.Sprintf("%s", presignedURL),
	})
}
