package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"takeaway-go/config"
	"takeaway-go/pkg/database"
	"takeaway-go/pkg/redis"
	"takeaway-go/router"
)

func main() {
	// 初始化配置文件
	config.Init()

	// 初始化数据库
	database.InitDB()
	// 初始化 Redis
	redis.InitRedis()

	// 初始化 *gin.Engine
	gin.SetMode(config.AppConf.Server.Mode)
	r := gin.Default()
	// 注册路由
	router.InitRouter(r)

	// 启动服务
	port := fmt.Sprintf(":%d", config.AppConf.Server.Port)
	if err := r.Run(port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
