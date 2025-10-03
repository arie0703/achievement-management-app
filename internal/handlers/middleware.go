package handlers

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse エラーレスポンス形式
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ValidationError バリデーションエラー
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}



// CORSMiddleware CORS設定ミドルウェア
func (s *Server) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}