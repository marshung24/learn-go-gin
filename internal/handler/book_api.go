package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookResponse 是 API 回應的書籍資料結構
// json tag 定義 JSON 欄位名稱（小寫駝峰式）
type BookResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	ISBN      string `json:"isbn"`
	Stock     int    `json:"stock"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// fakeBooksAPI 是 API 使用的假資料
// 後續 U06 會改為從資料庫取得
var fakeBooksAPI = []BookResponse{
	{ID: 1, Title: "Clean Code", Author: "Robert C. Martin", ISBN: "978-0132350884", Stock: 5, CreatedAt: "2024-01-01T10:00:00Z", UpdatedAt: "2024-01-01T10:00:00Z"},
	{ID: 2, Title: "The Pragmatic Programmer", Author: "Andy Hunt", ISBN: "978-0135957059", Stock: 0, CreatedAt: "2024-01-02T11:00:00Z", UpdatedAt: "2024-01-02T11:00:00Z"},
	{ID: 3, Title: "Refactoring", Author: "Martin Fowler", ISBN: "978-0134757599", Stock: 3, CreatedAt: "2024-01-03T12:00:00Z", UpdatedAt: "2024-01-03T12:00:00Z"},
}

// GetBooks godoc
// @Summary     取得所有書籍
// @Description 取得書籍清單
// @Tags        books
// @Produce     json
// @Success     200 {array} BookResponse
// @Router      /api/books [get]
func GetBooks(c *gin.Context) {
	// c.JSON() 自動設定 Content-Type: application/json
	c.JSON(http.StatusOK, fakeBooksAPI)
}

// GetBookByID godoc
// @Summary     依 ID 取得書籍
// @Description 依書籍 ID 取得單筆書籍資料
// @Tags        books
// @Produce     json
// @Param       id path int true "書籍 ID"
// @Success     200 {object} BookResponse
// @Failure     404 {object} map[string]string
// @Router      /api/books/{id} [get]
func GetBookByID(c *gin.Context) {
	// c.Param("id") 取得路徑參數
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 ID"})
		return
	}

	// 從假資料中找書（後續 U06 改為查 DB）
	for _, book := range fakeBooksAPI {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "找不到書籍"})
}

// CreateBookAPI godoc
// @Summary     新增書籍
// @Description 新增一本書籍
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       book body object true "書籍資料"
// @Success     201 {object} BookResponse
// @Failure     400 {object} map[string]string
// @Router      /api/books [post]
func CreateBookAPI(c *gin.Context) {
	// 目前只是假實作，印出收到的資料
	var input struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		ISBN   string `json:"isbn"`
		Stock  int    `json:"stock"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 假回應（後續 U06 會存入 DB 並回傳真實資料）
	response := BookResponse{
		ID:        4, // 假 ID
		Title:     input.Title,
		Author:    input.Author,
		ISBN:      input.ISBN,
		Stock:     input.Stock,
		CreatedAt: "2024-01-04T10:00:00Z",
		UpdatedAt: "2024-01-04T10:00:00Z",
	}

	// HTTP 201 Created
	c.JSON(http.StatusCreated, response)
}

// UpdateBookAPI godoc
// @Summary     更新書籍
// @Description 依 ID 更新書籍資料
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       id path int true "書籍 ID"
// @Param       book body object true "書籍資料"
// @Success     200 {object} BookResponse
// @Failure     404 {object} map[string]string
// @Router      /api/books/{id} [put]
func UpdateBookAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 ID"})
		return
	}

	var input struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		ISBN   string `json:"isbn"`
		Stock  int    `json:"stock"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 假回應（後續 U06 會真正更新 DB）
	response := BookResponse{
		ID:        id,
		Title:     input.Title,
		Author:    input.Author,
		ISBN:      input.ISBN,
		Stock:     input.Stock,
		CreatedAt: "2024-01-01T10:00:00Z",
		UpdatedAt: "2024-01-04T15:00:00Z",
	}

	c.JSON(http.StatusOK, response)
}

// DeleteBookAPI godoc
// @Summary     刪除書籍
// @Description 依 ID 刪除書籍
// @Tags        books
// @Param       id path int true "書籍 ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Router      /api/books/{id} [delete]
func DeleteBookAPI(c *gin.Context) {
	// 目前只是假實作，後續 U06 會真正從 DB 刪除
	// HTTP 204 No Content（刪除成功，無回應內容）
	c.Status(http.StatusNoContent)
}
