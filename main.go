package main

import (
	"github.com/example/learn-go-gin/internal/handler"
	"github.com/gin-gonic/gin"
)

// main 是 Go 應用程式的進入點。
// 使用 Gin 框架建立 HTTP 伺服器，並註冊路由。
// gin.Default() 會自動包含 Logger 和 Recovery middleware：
// - Logger：記錄每個請求的方法、路徑、狀態碼、處理時間
// - Recovery：捕獲 panic 並回傳 500 錯誤，避免程式崩潰
func main() {
	// 建立 Gin 引擎，包含預設的 Logger 和 Recovery middleware
	r := gin.Default()

	// 載入 HTML 模板
	// LoadHTMLGlob 會載入指定 pattern 的所有模板檔案
	r.LoadHTMLGlob("templates/**/*")

	// 基本路由
	r.GET("/hello", handler.Hello)
	r.GET("/health", handler.Health)
	r.GET("/whoami", handler.Whoami)

	// 書籍 MVC 路由（頁面渲染）
	r.GET("/books", handler.ListBooks)
	r.GET("/books/new", handler.NewBookForm)
	r.GET("/books/:id", handler.ShowBook)
	r.GET("/books/:id/edit", handler.EditBookForm)
	r.POST("/books", handler.CreateBook)
	r.POST("/books/:id/edit", handler.UpdateBook)
	r.POST("/books/:id/delete", handler.DeleteBook)

	// 書籍 REST API 路由（回傳 JSON）
	// r.Group() 建立路由群組，統一路徑前綴
	api := r.Group("/api")
	{
		api.GET("/books", handler.GetBooks)
		api.GET("/books/:id", handler.GetBookByID)
		api.POST("/books", handler.CreateBookAPI)
		api.PUT("/books/:id", handler.UpdateBookAPI)
		api.DELETE("/books/:id", handler.DeleteBookAPI)
	}

	// 啟動 HTTP 伺服器，監聽 8080 port
	// 這是內嵌的 HTTP Server，不需要額外安裝 Web Server
	r.Run(":8080")
}
