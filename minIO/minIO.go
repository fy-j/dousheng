package minIO

import (
	"dousheng/config"
	"fmt"
	_ "github.com/minio/minio-go/pkg/encrypt"
	"github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/policy"
	"log"
	"net/url"
	"time"
)

// 全局变量
var (
	Client *minio.Client
)

// InitClient : 连接 minIO 返回对应client
func InitMinIO(cfg *config.MinIOConfig) {
	// 初使化minio client对象。
	var err error
	if Client, err = minio.New(cfg.Endpoint, cfg.Username, cfg.Password, false); err != nil {
		log.Fatalln(err)
		return
	}
	//MinIO桶名称不能带下划线、只能小写字母
	CreateMinioBucket("userfeed")
}

//创建名称为bucketName 的视频流桶
func CreateMinioBucket(bucketName string) {
	location := "us-east-1"
	err := Client.MakeBucket(bucketName, location)
	if err != nil {
		exist, err := Client.BucketExists(bucketName)
		fmt.Println(exist)
		if err != nil && exist {
			fmt.Printf("We already own %s\n", bucketName)
		} else {
			fmt.Println(err, exist)
			return
		}
	}
	err = Client.SetBucketPolicy(bucketName, policy.BucketPolicyReadWrite)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Successfully created %s\n", bucketName)
}

func upload(bucketName string, fileName string, expires time.Duration) string {
	reqParams := make(url.Values)
	presignedURL, err := Client.PresignedGetObject(bucketName, fileName, expires, reqParams)
	if err != nil {
		//zap.L().Error(err.Error())
		return ""
	}
	return fmt.Sprintf("%s", presignedURL)
}
