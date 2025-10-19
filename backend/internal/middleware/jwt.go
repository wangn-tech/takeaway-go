package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"takeaway-go/common/enum"
	"takeaway-go/common/result"
	"takeaway-go/internal/app/config"
	"takeaway-go/internal/utils"
	"takeaway-go/pkg/logger"
)

// JWTAuth 创建一个 JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		authHeader := c.Request.Header.Get(config.AppConf.JWT.Admin.Name)
		if authHeader == "" {
			logger.Log.Warn("请求未携带 token，无权限访问", zap.String("path", c.Request.URL.Path))
			result.Fail(c, http.StatusUnauthorized, "请求未携带token，无权限访问")
			c.Abort()
			return
		}

		// --- 这里是修改点 ---
		// 直接将 authHeader 作为 tokenString 进行解析
		// 不再检查 "Bearer " 前缀
		claims, err := utils.ParseToken(authHeader, config.AppConf.JWT.Admin.Secret)
		if err != nil {
			// token 解析失败
			logger.Log.Warn("token 解析失败", zap.Error(err), zap.String("path", c.Request.URL.Path))
			// 根据不同的错误类型返回不同的消息
			if strings.Contains(err.Error(), "token is expired") {
				result.Fail(c, http.StatusUnauthorized, "token已过期")
			} else {
				result.Fail(c, http.StatusUnauthorized, "无效的token")
			}
			c.Abort()
			return
		}

		// 将当前请求的 claims 信息保存到请求的上下文 c 上
		// 后续的处理函数可以用c.Get("claims")来获取当前请求的用户信息
		c.Set("claims", claims)
		c.Set(enum.CurrentId, claims.ID)
		c.Set(enum.CurrentName, claims.GrantScope)
		c.Next()
	}
}
