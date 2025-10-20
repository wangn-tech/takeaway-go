package router

import (
	"takeaway-go/internal/api/handler"
	"takeaway-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化并注册所有路由
func InitRouter(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 初始化依赖
	// employeeRepo := repository.NewEmployeeRepository(database.DB)
	// employeeService := service.NewEmployeeService(employeeRepo)
	// employeeHandler := handler.NewEmployeeHandler(employeeService)

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
			// 修改员工状态
			adminRouter.POST("/employee/status/:status", employeeHandler.UpdateStatus)
			// 修改密码
			adminRouter.PUT("/employee/editPassword", employeeHandler.EditPassword)
			// 更改用户信息
			adminRouter.PUT("/employee", employeeHandler.UpdateEmployee)
			// 员工分页查询
			adminRouter.GET("/employee/page", employeeHandler.PageQuery)
			// 根据ID获取员工信息
			adminRouter.GET("/employee/:id", employeeHandler.GetByID)
		}
	}
}
