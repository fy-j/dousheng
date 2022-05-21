package minIO

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v6"
)

func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	fileObj, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err := Client.PutObject("userfeed", "shipin1", fileObj, file.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		c.JSON(0, gin.H{ // H是一个开箱即用的map
			"message": "fill",
		})
	}
	fmt.Println("Successfully uploaded bytes: ", n)
	c.JSON(200, gin.H{ // H是一个开箱即用的map
		"message": "success",
	})
}
