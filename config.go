package main

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf Config

type MySQLConfig struct {
	Addr string
}

type RedisConfig struct {
	Addr string
}

// 读取配置文件config
type Config struct {
	Redis      RedisConfig
	MySQL      MySQLConfig
	ListenPort int
}

func ConfInit() {
	viper.SetConfigName("config.toml") // 配置文件名 (不带扩展格式)
	viper.AddConfigPath("./")          // 配置文件的路径
	err := viper.ReadInConfig()        //找到并读取配置文件
	if err != nil {                    // 捕获读取中遇到的error
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	viper.Unmarshal(&Conf) //将配置文件绑定到config上
	fmt.Println("Conf: ", Conf)
}
