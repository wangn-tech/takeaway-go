package handler

import (
	"net/http"
	"strconv"
	"takeaway-go/common/enum"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/service"
	"takeaway-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmployeeHandler struct {
	// service *service.EmployeeService
	service service.IEmployeeService
}

func NewEmployeeHandler(service service.IEmployeeService) *EmployeeHandler {
	// return &EmployeeHandler{
	// 	service: service.NewEmployeeService(),
	// }
	return &EmployeeHandler{service: service}
}

// Login 处理员工登录请求
func (h *EmployeeHandler) Login(ctx *gin.Context) {
	var loginDTO request.EmployeeLoginDTO
	// 绑定和校验请求参数
	if err := ctx.ShouldBindJSON(&loginDTO); err != nil {
		logger.Log.Warn("Login: 参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "参数无效")
		return
	}

	// 调用 service 进行登录
	loginVO, err := h.service.Login(ctx.Request.Context(), loginDTO)
	if err != nil {
		logger.Log.Warn("Login: 登录失败", zap.String("username", loginDTO.Username), zap.Error(err))
		result.Fail(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	// 登录成功，返回令牌和用户信息
	result.Success(ctx, "登录成功", loginVO)
}

// Logout 处理员工登出请求
//
//	// 在 JWT 方案中, 服务端的登出是可选的, 真正的登出由客户端删除 token 实现
//	// 此处实现: 在 redis 中删除 token, 使其失效
func (h *EmployeeHandler) Logout(ctx *gin.Context) {
	// 从 Gin 上下文中获取由 JWT 中间件设置的 claims
	// claimsValue, exists := ctx.Get("claims")
	userID, exists := ctx.Get(enum.CurrentId)
	if !exists {
		logger.Log.Error("Logout: 获取用户ID失败")
		result.Fail(ctx, http.StatusUnauthorized, "无法获取用户信息")
		return
	}

	// 调用 Service 层处理登出逻辑
	if err := h.service.Logout(ctx.Request.Context(), userID.(uint64)); err != nil {
		logger.Log.Error("Logout: 登出失败", zap.Uint64("userId", userID.(uint64)), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回成功响应
	result.Success(ctx, "登出成功", nil)
}

// EditPassword 处理修改员工密码请求
func (h *EmployeeHandler) EditPassword(ctx *gin.Context) {
	var req request.EmployeeEditPasswordDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("EditPassword: 参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// 从 context 中获取用户 id
	if id, ok := ctx.Get(enum.CurrentId); ok {
		req.EmpId = id.(uint64)
	}
	err := h.service.EditPassword(ctx.Request.Context(), req)
	if err != nil {
		logger.Log.Warn("EditPassword: 修改密码失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回成功响应
	result.Success(ctx, "修改密码成功，请重新登录", nil)
}

// AddEmployee 处理新增员工请求
func (h *EmployeeHandler) AddEmployee(ctx *gin.Context) {
	var addDTO request.EmployeeAddDTO
	// 绑定和校验参数
	if err := ctx.ShouldBindJSON(&addDTO); err != nil {
		logger.Log.Warn("AddEmployee: 参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, "参数无效")
		return
	}

	// 调用 Service 层处理业务逻辑
	if err := h.service.AddEmployee(ctx.Request.Context(), addDTO); err != nil {
		logger.Log.Warn("AddEmployee: 新增员工失败", zap.Any("dto", addDTO), zap.Error(err))
		// 如果是业务错误（如用户名已存在），返回 200 和错误信息
		result.Fail(ctx, http.StatusOK, err.Error())
		return
	}

	// 3. 返回成功响应
	result.Success(ctx, "新增员工成功", nil)
}

// PageQuery 处理员工分页查询请求
func (h *EmployeeHandler) PageQuery(ctx *gin.Context) {
	var pageQueryDTO request.EmployeePageQueryDTO

	// 绑定查询参数 (form)，如果绑定失败则使用默认值
	if err := ctx.ShouldBindQuery(&pageQueryDTO); err != nil {
		// 对于分页查询，即使参数有误，也应提供默认查询结果，而不是直接报错
		logger.Log.Warn("PageQuery: 查询参数绑定失败，使用默认值", zap.Error(err))
		pageQueryDTO.Page = 1
		pageQueryDTO.PageSize = 10
	}

	// 调用 Service 层处理业务逻辑
	pageResult, err := h.service.PageQuery(ctx.Request.Context(), pageQueryDTO)
	if err != nil {
		logger.Log.Error("PageQuery: 分页查询失败", zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	// 返回成功响应
	result.Success(ctx, "查询成功", pageResult)
}

// Update 编辑员工信息
func (h *EmployeeHandler) UpdateEmployee(ctx *gin.Context) {
	var req request.EmployeeEditDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Update: 参数绑定失败", zap.Error(err))
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.UpdateEmployee(ctx.Request.Context(), &req)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, "编辑员工信息成功", nil)
}

func (h *EmployeeHandler) UpdateStatus(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Query("id"), 10, 64)
	status, _ := strconv.Atoi(ctx.Param("status"))

	err := h.service.UpdateStatus(ctx.Request.Context(), id, status)
	if err != nil {
		logger.Log.Warn("UpdateStatus: 更新员工状态失败", zap.Uint64("id", id), zap.Int("status", status), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, "更新员工状态成功", nil)
}

// GetByID 获取员工信息根据id
func (h *EmployeeHandler) GetByID(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	employee, err := h.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		logger.Log.Warn("GetById: 获取员工信息失败", zap.Uint64("id", id), zap.Error(err))
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, "获取员工信息成功", employee)
}
