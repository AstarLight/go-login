package main

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf Config

type MySQLConfig struct {
	Addr     string
	MaxIdle  int
	MaxConn  int
	ShowSQL  bool
	LogLevel int
}

type RedisConfig struct {
	Addr string
}

type SessionConfig struct {
	Key string
	TTL int
}

type SqlConfig struct {
	File string
}

type BlacklistConfig struct {
	IPList       []string
	UsernameList []string
}

type LimitCountConfig struct {
	Limit     int
	Mode      string
	KeyPrefix string
}

type CommonConfig struct {
	ListenPort         int
	EnableRedisCache   bool
	EnableLocalCache   bool
	MinPasswordLength  int
	MaxUsernameLen     int
	EnterPage          string
	HomePage           string
	LoginPage          string
	RegistPage         string
	PasswordComplexity []string
}

// 读取配置文件config
type Config struct {
	Redis     RedisConfig
	MySQL     MySQLConfig
	SQL       SqlConfig
	Session   SessionConfig
	Common    CommonConfig
	Blacklist BlacklistConfig
	RateLimit LimitCountConfig
}

func ConfInit() {
	viper.SetConfigFile("config.toml") // 配置文件名 (不带扩展格式)
	viper.AddConfigPath("./")          // 配置文件的路径
	err := viper.ReadInConfig()        //找到并读取配置文件
	if err != nil {                    // 捕获读取中遇到的error
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	viper.Unmarshal(&Conf) //将配置文件绑定到config上
	fmt.Printf("Conf: %+v\n", Conf)
}
