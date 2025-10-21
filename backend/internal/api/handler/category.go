package handler

import (
	"net/http"
	"strconv"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/service"
	"takeaway-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	service service.ICategoryService
}

func NewCategoryHandler(service service.ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

// AddCategory 新增分类
func (h *CategoryHandler) AddCategory(ctx *gin.Context) {
	var dto request.CategoryDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logger.Log.Warn("Add: 参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AddCategory(ctx.Request.Context(), dto); err != nil {
		logger.Log.Warn("Add: 新增分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, "新增分类成功", nil)
}

// PageQuery 分页查询分类
func (h *CategoryHandler) PageQuery(ctx *gin.Context) {
	var req request.CategoryPageQueryDTO
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Log.Warn("PageQuery: 参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	pageResult, err := h.service.PageQuery(ctx.Request.Context(), req)
	if err != nil {
		logger.Log.Warn("PageQuery: 分页查询分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, "查询成功", pageResult)
}

// List 查询分类列表
func (h *CategoryHandler) List(ctx *gin.Context) {
	// 获取参数
	cate, err := strconv.Atoi(ctx.Query("type"))
	if err != nil {
		logger.Log.Warn("List: 获取分类类型失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "分类类型参数缺失")
		return
	}

	// 查询
	categories, err := h.service.List(ctx.Request.Context(), cate)
	if err != nil {
		logger.Log.Warn("List: 查询分类列表失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, "查询成功", categories)
}

// DeleteById 删除分类
func (h *CategoryHandler) DeleteById(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Fail(ctx, http.StatusBadRequest, "分类ID不能为空")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		result.Fail(ctx, http.StatusBadRequest, "分类ID格式错误")
		return
	}

	if err := h.service.DeleteById(ctx.Request.Context(), id); err != nil {
		logger.Log.Warn("Delete: 删除分类失败", zap.Uint64("id", id), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, "删除分类成功", nil)
}

// UpdateStatus 启用/禁用分类
func (h *CategoryHandler) UpdateStatus(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Query("id"), 10, 64)
	status, _ := strconv.Atoi(ctx.Param("status"))
	err := h.service.SetStatus(ctx, id, status)
	if err != nil {
		logger.Log.Warn("UpdateStatus: 更新分类状态失败", zap.Uint64("id", id), zap.Int("status", status), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, "更新分类状态成功", nil)
}

// EditCategory 修改分类
func (h *CategoryHandler) EditCategory(ctx *gin.Context) {
	var dto request.CategoryDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		logger.Log.Warn("EditCategory: 参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Update(ctx.Request.Context(), dto); err != nil {
		logger.Log.Warn("EditCategory: 更新分类失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, "更新分类成功", nil)
}
