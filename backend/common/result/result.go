package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"takeaway-go/common/e"
)

func Response(c *gin.Context, httpStatus int, code int, msg string, data interface{}) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// Success 成功响应
func Success(c *gin.Context, msg string, data interface{}) {
	Response(c, http.StatusOK, e.SUCCESS, msg, data)
}

// Fail 失败响应: 业务错误, 非系统错误
func Fail(c *gin.Context, code int, msg string) {
	Response(c, http.StatusOK, code, msg, nil)
}

// Fatal 系统错误响应
func Fatal(c *gin.Context, httpStatus int, code int, msg string) {
	Response(c, httpStatus, code, msg, nil)
}
