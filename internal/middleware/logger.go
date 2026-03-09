package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 自訂的 request logging middleware
// 記錄每個請求的方法、路徑、狀態碼、處理時間
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 記錄開始時間
		start := time.Now()

		// 處理請求
		c.Next()

		// 計算處理時間
		duration := time.Since(start)

		// 記錄請求資訊
		log.Printf("[%s] %s %d %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
