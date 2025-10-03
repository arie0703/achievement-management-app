package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandlerMiddleware(t *testing.T) {
	// テスト用のGinエンジンを作成
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandlerMiddleware())

	// エラーを発生させるテストハンドラー
	router.GET("/test-error", func(c *gin.Context) {
		_ = c.Error(assert.AnError).SetType(gin.ErrorTypePublic)
	})

	// バインドエラーを発生させるテストハンドラー
	router.POST("/test-bind-error", func(c *gin.Context) {
		_ = c.Error(assert.AnError).SetType(gin.ErrorTypeBind)
	})

	// 内部エラーを発生させるテストハンドラー
	router.GET("/test-internal-error", func(c *gin.Context) {
		_ = c.Error(assert.AnError).SetType(gin.ErrorTypePrivate)
	})

	// 正常なレスポンスのテストハンドラー
	router.GET("/test-success", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Public Error",
			method:         "GET",
			path:           "/test-error",
			expectedStatus: 400,
			expectedError:  "bad_request",
		},
		{
			name:           "Bind Error",
			method:         "POST",
			path:           "/test-bind-error",
			expectedStatus: 400,
			expectedError:  "validation_error",
		},
		{
			name:           "Internal Error",
			method:         "GET",
			path:           "/test-internal-error",
			expectedStatus: 500,
			expectedError:  "internal_error",
		},
		{
			name:           "Success Response",
			method:         "GET",
			path:           "/test-success",
			expectedStatus: 200,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			}
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	// テスト用のGinエンジンを作成
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())

	// テストハンドラー
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// OPTIONSリクエストのテスト
	t.Run("OPTIONS Request", func(t *testing.T) {
		req, err := http.NewRequest("OPTIONS", "/test", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	})

	// 通常のリクエストのテスト
	t.Run("Normal Request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	})
}