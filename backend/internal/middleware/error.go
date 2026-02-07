package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果已经响应，则不再处理
		if c.Writer.Written() {
			return
		}

		// 处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 根据错误类型返回不同的状态码
			statusCode := http.StatusInternalServerError
			message := "internal server error"

			switch err {
			case http.ErrHandlerTimeout:
				statusCode = http.StatusRequestTimeout
				message = "request timeout"
			default:
				if _, ok := err.(interface{ Timeout() bool }); ok {
					statusCode = http.StatusRequestTimeout
					message = "request timeout"
				}
			}

			c.JSON(statusCode, ErrorResponse{
				Success: false,
				Message: message,
				Code:    statusCode,
			})
		}
	}
}

// RecoveryHandler 自定义恢复中间件
func RecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Success: false,
					Message: "internal server error",
					Code:    http.StatusInternalServerError,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
