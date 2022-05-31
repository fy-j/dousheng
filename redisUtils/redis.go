package redisUtils

import (
	"dousheng/config"
	"dousheng/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var (
	Clients *redis.Client
	Nil     = redis.Nil
)

type SliceCmd = redis.SliceCmd
type StringStringMapCmd = redis.StringStringMapCmd

// Init 初始化连接
func Init(cfg *config.RedisConfig) (err error) {
	Clients = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	_, err = Clients.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = Clients.Close()
}

//从redis中获取发布视频切片
func GetVideoInfoListFromRedis(key_type string, userId int) (*[]model.VideoInfo, error) {
	key := Generate(key_type, strconv.FormatInt(int64(userId), 10))
	result, err := Clients.Get(key).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	} else {
		var info []model.VideoInfo
		if err := json.Unmarshal([]byte(result), &info); err != nil {
			return nil, err
		}
		return &info, nil
	}
}

func Set(key string, any interface{}, time time.Duration) error {
	jsonstring, err := json.Marshal(any)
	if err != nil {
		return err
	} else {
		errSet := Clients.Set(key, string(jsonstring), time).Err()
		if errSet != nil {
			return errSet
		} else {
			return nil
		}
	}
}
