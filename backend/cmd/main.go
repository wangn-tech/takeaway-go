package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动服务，默认监听在 :8080 端口
	err := r.Run(":8082")
	if err != nil {
		return
	}
}
