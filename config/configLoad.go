package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(Config)

//
type Config struct {
	Port        int `mapstructure:"port"`
	Redis		RedisConfig `mapstructure:"redis"`
	MinIO		MinIOConfig `mapstructure:"minio"`
	Mongo		MongoConfig `mapstructure:"mongo"`
}

//redis 配置类
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

//minIO配置类
type MinIOConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

//mongodb config struct
type MongoConfig struct {
	DBname string ` json:"name" `
	Host   string ` json:"host" `
	Port   string ` json:"port" `
	User   string ` json:"user" `
	Pwd    string ` json:"pwd"  `
}

func Init() error {
	viper.SetConfigFile("config/config.yaml")

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改...")
		viper.Unmarshal(&Conf)
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unmarshal to Conf failed, err:%v", err))
	}
	return err
}
