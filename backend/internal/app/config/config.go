package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

// AppConf 是全局配置变量
var AppConf *Config

// Config 是应用程序的主配置结构体
type Config struct {
	Server ServerConfig   `mapstructure:"server"`
	MySQL  DatabaseConfig `mapstructure:"mysql"`
	Redis  RedisConfig    `mapstructure:"redis"`
	JWT    JWTConfig      `mapstructure:"jwt"`
	Log    LogConfig      `mapstructure:"log"`
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

// JWTConfig JWT 配置
type JWTConfig struct {
	Admin JWTOption `mapstructure:"admin"`
	User  JWTOption `mapstructure:"user"`
}

// JWTOption JWT 选项
type JWTOption struct {
	Secret string `mapstructure:"secret"`
	Name   string `mapstructure:"name"`
}

// LogConfig 定义了日志的配置参数
type LogConfig struct {
	Level      string `mapstructure:"level"`      // 日志级别, 例如: debug, info, warn, error
	Format     string `mapstructure:"format"`     // 日志格式, 例如: console, json
	ShowLine   bool   `mapstructure:"show-line"`  // 是否显示行号
	Stacktrace bool   `mapstructure:"stacktrace"` // 是否开启堆栈跟踪
}

// Init 读取解析配置文件
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
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 解析配置
	AppConf = &Config{}
	if err := config.Unmarshal(AppConf); err != nil {
		panic(fmt.Errorf("unable to decode config: %w", err))
	}

	log.Println("Configuration loaded successfully")
}
