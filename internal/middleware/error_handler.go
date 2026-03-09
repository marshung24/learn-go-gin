package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppError 自訂錯誤類型，用於統一錯誤格式
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 實作 error 介面
func (e *AppError) Error() string {
	return e.Message
}

// 常用錯誤定義
var (
	ErrNotFound      = &AppError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrBadRequest    = &AppError{Code: http.StatusBadRequest, Message: "Invalid request"}
	ErrDuplicateISBN = &AppError{Code: http.StatusConflict, Message: "ISBN already exists"}
	ErrInternal      = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
)

// ErrorResponse 錯誤回應格式
type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Path   string `json:"path"`
}

// ErrorHandler 全域錯誤處理 Middleware
// 集中處理所有 error，回傳統一格式的 JSON 錯誤回應
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 檢查是否有錯誤
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 處理 GORM 的 ErrRecordNotFound
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, ErrorResponse{
					Status: http.StatusNotFound,
					Error:  "Not Found",
					Path:   c.Request.URL.Path,
				})
				return
			}

			// 處理自訂 AppError
			if appErr, ok := err.(*AppError); ok {
				c.JSON(appErr.Code, ErrorResponse{
					Status: appErr.Code,
					Error:  appErr.Message,
					Path:   c.Request.URL.Path,
				})
				return
			}

			// 未知錯誤，記錄 log 並回傳 500
			log.Printf("Unexpected error on %s: %v", c.Request.URL.Path, err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Status: http.StatusInternalServerError,
				Error:  "Internal Server Error",
				Path:   c.Request.URL.Path,
			})
		}
	}
}

// Recovery 自訂的 panic 恢復 middleware
// 攔截 panic 避免整個 server 崩潰
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered on %s: %v", c.Request.URL.Path, err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Status: http.StatusInternalServerError,
					Error:  "Internal Server Error",
					Path:   c.Request.URL.Path,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
