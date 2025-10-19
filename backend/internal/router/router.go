package router

import (
	"github.com/gin-gonic/gin"
	"takeaway-go/internal/api/handler"
	"takeaway-go/internal/middleware"
)

// InitRouter 初始化并注册所有路由
func InitRouter(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	employeeHandler := handler.NewEmployeeHandler()

	adminRouter := r.Group("/admin")
	{
		adminRouter.POST("/employee/login", employeeHandler.Login)
		// 对该路由组下的其他路由使用JWT认证中间件
		adminRouter.Use(middleware.JWTAuth())
		{
			// 员工登出
			adminRouter.POST("/employee/logout", employeeHandler.Logout)
			// 新增员工
			adminRouter.POST("/employee", employeeHandler.AddEmployee)
			// 员工分页查询
			adminRouter.GET("/employee/page", employeeHandler.PageQuery)
		}
	}
}
