package admin

import (
	"takeaway-go/internal/api/handler"
	"takeaway-go/internal/app/initializer/database"
	"takeaway-go/internal/middleware"
	"takeaway-go/internal/repository"
	"takeaway-go/internal/service"

	"github.com/gin-gonic/gin"
)

type EmployeeRouter struct {
	service service.IEmployeeService
}

func (r *EmployeeRouter) InitApiRouter(router *gin.RouterGroup) {
	// 依赖注入
	r.service = service.NewEmployeeService(
		repository.NewEmployeeDao(database.DB),
	)

	employeeHandler := handler.NewEmployeeHandler(r.service)

	// "/employee" 路由组
	employeeGroup := router.Group("/employee")
	{
		// 员工登录
		employeeGroup.POST("/login", employeeHandler.Login)
	}
	// 使用 JWTAuth 中间件
	employeeGroup.Use(middleware.JWTAuth())
	{
		// 退出登录
		employeeGroup.POST("/logout", employeeHandler.Logout)
		// 新增员工
		employeeGroup.POST("", employeeHandler.AddEmployee)
		// 员工分页查询
		employeeGroup.GET("/page", employeeHandler.PageQuery)
		// 启用、禁用员工
		employeeGroup.GET("/status/:status", employeeHandler.UpdateStatus)
		// 编辑员工信息
		employeeGroup.PUT("", employeeHandler.UpdateEmployee)
		// 根据 id 查询员工
		employeeGroup.GET("/:id", employeeHandler.GetByID)
		// 修改密码
		employeeGroup.PUT("/editPassword", employeeHandler.EditPassword)
	}
}
