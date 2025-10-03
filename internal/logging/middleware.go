package logging

import (
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware HTTPリクエストのログを記録するミドルウェア
func LoggingMiddleware(accessLogger *AccessLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// リクエストを処理
		c.Next()
		
		// ログを記録
		duration := time.Since(start)
		accessLogger.LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			duration,
		)
	}
}

// ErrorLoggingMiddleware エラーログを記録するミドルウェア
func ErrorLoggingMiddleware(errorLogger *ErrorLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// エラーがある場合はログに記録
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				errorLogger.LogAPIError(
					c.Request.URL.Path,
					c.Request.Method,
					c.Writer.Status(),
					err.Err,
				)
			}
		}
	}
}

// RecoveryMiddleware パニックからの回復とログ記録
func RecoveryMiddleware(errorLogger *ErrorLogger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(error); ok {
			errorLogger.LogError("panic_recovery", "middleware", err, map[string]interface{}{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			})
		}
		c.AbortWithStatus(500)
	})
}