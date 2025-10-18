package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

// AppConf 是全局配置变量
var AppConf *Config

type Config struct {
	Server ServerConfig   `mapstructure:"server"`
	MySQL  DatabaseConfig `mapstructure:"mysql"`
	Redis  RedisConfig    `mapstructure:"redis"`
}

// ServerConfig 后端服务配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库 MySQL 配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Init() {
	// 解析命令行参数 “env”, 默认为 "dev"
	env := pflag.String("env", "dev", "Specify the environment config to use: [dev, prod, test]")
	pflag.Parse()

	// 设置配置文件: ./config/config-{env}.yaml
	config := viper.New()
	config.AddConfigPath("./config")
	config.SetConfigName(fmt.Sprintf("config-%s", *env))
	config.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 解析配置
	AppConf = &Config{}
	if err := config.Unmarshal(AppConf); err != nil {
		panic(fmt.Errorf("unable to decode config: %w", err))
	}

	log.Println("Configuration loaded successfully")
}
