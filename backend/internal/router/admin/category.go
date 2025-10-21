package admin

import (
	"takeaway-go/internal/api/handler"
	"takeaway-go/internal/app/initializer/database"
	"takeaway-go/internal/middleware"
	"takeaway-go/internal/repository"
	"takeaway-go/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryRouter struct {
	service service.ICategoryService
}

func (r *CategoryRouter) InitApiRouter(router *gin.RouterGroup) {
	// 依赖注入
	r.service = service.NewCategoryService(
		repository.NewCategoryDao(database.DB),
	)
	categoryHandler := handler.NewCategoryHandler(r.service)

	categoryGroup := router.Group("/category")
	// 使用 JWTAuth 中间件
	categoryGroup.Use(middleware.JWTAuth())
	{
		categoryGroup.POST("", categoryHandler.AddCategory)                 // 新增分类
		categoryGroup.GET("/page", categoryHandler.PageQuery)               // 分类分页查询
		categoryGroup.GET("/list", categoryHandler.List)                    // 根据类型查询分类 (查询分类列表)
		categoryGroup.DELETE("", categoryHandler.DeleteById)                // 根据 id 删除分类
		categoryGroup.PUT("", categoryHandler.EditCategory)                 // 修改分类
		categoryGroup.POST("/status/:status", categoryHandler.UpdateStatus) // 启用/禁用分类
	}
}
