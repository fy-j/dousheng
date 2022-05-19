package minIO

import (
	"dousheng/config"
	_ "github.com/minio/minio-go/pkg/encrypt"
	"github.com/minio/minio-go/v6"
	"log"
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
}
