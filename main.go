package main

import (
	"github.com/example/learn-go-gin/internal/config"
	"github.com/example/learn-go-gin/internal/handler"
	"github.com/example/learn-go-gin/internal/middleware"
	"github.com/example/learn-go-gin/internal/repository"
	"github.com/example/learn-go-gin/internal/service"
	"github.com/gin-gonic/gin"
)

// main 是 Go 應用程式的進入點。
// 使用 Gin 框架建立 HTTP 伺服器，並註冊路由。
// 分層架構：Handler → Service → Repository → DB
func main() {
	// 載入設定檔
	config.LoadConfig()

	// 初始化資料庫連線
	config.InitDB()

	// 初始化各層
	// Repository → Service → Handler
	bookRepo := repository.NewBookRepository(config.DB)
	bookService := service.NewBookService(bookRepo)
	bookViewHandler := handler.NewBookViewHandler(bookService)
	bookAPIHandler := handler.NewBookAPIHandler(bookService)

	// 建立 Gin 引擎
	// 使用 gin.New() 搭配自訂 middleware，而非 gin.Default()
	r := gin.New()

	// 註冊 middleware
	// - Recovery: 攔截 panic，避免 server 崩潰
	// - Logger: 記錄請求日誌
	// - ErrorHandler: 統一錯誤回應格式
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	// 載入 HTML 模板
	r.LoadHTMLGlob("templates/**/*")

	// 基本路由
	r.GET("/hello", handler.Hello)
	r.GET("/health", handler.Health)
	r.GET("/whoami", handler.Whoami)

	// 書籍 MVC 路由（頁面渲染）
	r.GET("/books", bookViewHandler.ListBooks)
	r.GET("/books/new", bookViewHandler.NewBookForm)
	r.GET("/books/:id", bookViewHandler.ShowBook)
	r.GET("/books/:id/edit", bookViewHandler.EditBookForm)
	r.POST("/books", bookViewHandler.CreateBook)
	r.POST("/books/:id/edit", bookViewHandler.UpdateBook)
	r.POST("/books/:id/delete", bookViewHandler.DeleteBook)

	// 書籍 REST API 路由（回傳 JSON）
	api := r.Group("/api")
	{
		api.GET("/books", bookAPIHandler.GetBooks)
		api.GET("/books/:id", bookAPIHandler.GetBookByID)
		api.POST("/books", bookAPIHandler.CreateBook)
		api.PUT("/books/:id", bookAPIHandler.UpdateBook)
		api.DELETE("/books/:id", bookAPIHandler.DeleteBook)
	}

	// 啟動 HTTP 伺服器
	r.Run(":" + config.AppCfg.App.Port)
}
