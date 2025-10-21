package router

import (
	"takeaway-go/internal/router/admin"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	admin.EmployeeRouter
	admin.CategoryRouter
}

var AllRouter = new(RouterGroup)

// InitRouter 初始化并注册所有路由
func InitRouter(r *gin.Engine) {

	allRouter := AllRouter

	// 测试路由ping, pong
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// admin
	adminRouter := r.Group("/admin")
	{
		allRouter.EmployeeRouter.InitApiRouter(adminRouter)
		allRouter.CategoryRouter.InitApiRouter(adminRouter)
	}

}
