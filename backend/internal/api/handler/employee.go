package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"takeaway-go/common/result"
	"takeaway-go/internal/api/request"
	"takeaway-go/internal/service"
	"takeaway-go/internal/utils"
	"takeaway-go/pkg/logger"
)

type EmployeeHandler struct {
	service *service.EmployeeService
}

func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{
		service: service.NewEmployeeService(),
	}
}

// Login 处理员工登录请求
func (h *EmployeeHandler) Login(c *gin.Context) {
	var loginDTO request.EmployeeLoginDTO
	// 绑定和校验请求参数
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		logger.Log.Warn("Login: 参数绑定失败", zap.Error(err))
		result.Fail(c, http.StatusBadRequest, "参数无效")
		return
	}

	// 调用 service 进行登录
	loginVO, err := h.service.Login(loginDTO)
	if err != nil {
		logger.Log.Warn("Login: 登录失败", zap.String("username", loginDTO.Username), zap.Error(err))
		result.Fail(c, http.StatusUnauthorized, err.Error())
		return
	}

	// 登录成功，返回令牌和用户信息
	result.Success(c, "登录成功", loginVO)
}

// Logout 处理员工登出请求
// 在 JWT 方案中, 服务端的登出是可选的, 真正的登出由客户端删除 token 实现
// 此处实现: 在 redis 中删除 token, 使其失效
func (h *EmployeeHandler) Logout(c *gin.Context) {
	// 从 Gin 上下文中获取由 JWT 中间件设置的 claims
	claimsValue, exists := c.Get("claims")
	if !exists {
		logger.Log.Error("Logout: 从上下文获取 claims 失败")
		result.Fail(c, http.StatusUnauthorized, "无法获取用户信息，请重新登录")
		return
	}
	// 类型断言，将 claims 转换为 *utils.Claims 类型
	claims, ok := claimsValue.(*utils.Claims)
	if !ok {
		logger.Log.Error("Logout: claims类型断言失败")
		result.Fail(c, http.StatusInternalServerError, "服务器内部错误，无法解析用户信息")
		return
	}

	// 调用 RevokeToken 函数，从 Redis 中删除用户的 token
	if err := utils.RevokeToken(claims.UserID); err != nil {
		logger.Log.Error("Logout: 撤销 token 失败", zap.Uint64("userId", claims.UserID), zap.Error(err))
		result.Fail(c, http.StatusInternalServerError, "登出失败")
		return
	}

	// 返回成功响应
	result.Success(c, "登出成功", nil)
}

// AddEmployee 处理新增员工请求
func (h *EmployeeHandler) AddEmployee(c *gin.Context) {
	var addDTO request.EmployeeAddDTO
	// 绑定和校验参数
	if err := c.ShouldBindJSON(&addDTO); err != nil {
		logger.Log.Warn("AddEmployee: 参数绑定失败", zap.Error(err))
		result.Fail(c, http.StatusBadRequest, "参数无效")
		return
	}

	// 调用 Service 层处理业务逻辑
	if err := h.service.AddEmployee(addDTO); err != nil {
		logger.Log.Warn("AddEmployee: 新增员工失败", zap.Any("dto", addDTO), zap.Error(err))
		// 如果是业务错误（如用户名已存在），返回 200 和错误信息
		result.Fail(c, http.StatusOK, err.Error())
		return
	}

	// 3. 返回成功响应
	result.Success(c, "新增员工成功", nil)
}

// PageQuery 处理员工分页查询请求
func (h *EmployeeHandler) PageQuery(c *gin.Context) {
	var pageQueryDTO request.EmployeePageQueryDTO

	// 绑定查询参数 (form)，如果绑定失败则使用默认值
	if err := c.ShouldBindQuery(&pageQueryDTO); err != nil {
		// 对于分页查询，即使参数有误，也应提供默认查询结果，而不是直接报错
		logger.Log.Warn("PageQuery: 查询参数绑定失败，使用默认值", zap.Error(err))
		pageQueryDTO.Page = 1
		pageQueryDTO.PageSize = 10
	}

	// 调用 Service 层处理业务逻辑
	pageResult, err := h.service.PageQuery(pageQueryDTO)
	if err != nil {
		logger.Log.Error("PageQuery: 分页查询失败", zap.Error(err))
		result.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	// 返回成功响应
	result.Success(c, "查询成功", pageResult)
}
