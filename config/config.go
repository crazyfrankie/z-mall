package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/kr/pretty"
	"github.com/spf13/viper"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env    string
	Server Server `yaml:"server"`
	MySQL  MySQL  `yaml:"mysql"`
	Redis  Redis  `yaml:"redis"`
	WeChat WeChat `yaml:"weChat"`
}

type Server struct {
	Host string `yaml:"host"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address string `yaml:"address"`
}

type WeChat struct {
	AppId  string `yaml:"appId"`
	AppKey string `yaml:"appKey"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	prefix := "config"
	contentFilePath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	viper.SetConfigFile(contentFilePath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	conf = new(Config)
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Failed to unmarshal key: %v", err)
	}

	conf.Env = GetEnv()
	// 打印配置，方便调试
	pretty.Printf("%+v\n", conf)
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}
