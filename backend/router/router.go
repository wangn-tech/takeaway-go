package router

import "github.com/gin-gonic/gin"

// InitRouter 初始化并注册所有路由
func InitRouter(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
