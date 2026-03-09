package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Hello 是最簡易的 handler，用於快速驗證應用程式是否正常運作。
// c.String() 會直接回傳純文字，適合簡單的測試端點。
func Hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, Gin!")
}

// Health 健康檢查端點，部署時可讓 load balancer / K8s probe 呼叫。
// 回傳 "OK" 表示應用程式正常運作。
func Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// Whoami 回傳使用者資訊（JSON 格式）。
// gin.H 是 map[string]any 的別名，方便建立 JSON 回應。
// c.JSON() 會自動設定 Content-Type: application/json。
func Whoami(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": "Go Gin Learner",
	})
}
