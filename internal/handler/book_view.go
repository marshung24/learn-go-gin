package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Book 代表書籍資料結構（目前使用假資料，後續 U06 會改為從 DB 取得）
type Book struct {
	ID        int64
	Title     string
	Author    string
	ISBN      string
	Stock     int
	CreatedAt string
	UpdatedAt string
}

// fakeBooks 是假資料，用於 U02 教學示範
// 後續 U06 會改為從資料庫取得
var fakeBooks = []Book{
	{ID: 1, Title: "Clean Code", Author: "Robert C. Martin", ISBN: "978-0132350884", Stock: 5, CreatedAt: "2024-01-01 10:00", UpdatedAt: "2024-01-01 10:00"},
	{ID: 2, Title: "The Pragmatic Programmer", Author: "Andy Hunt", ISBN: "978-0135957059", Stock: 0, CreatedAt: "2024-01-02 11:00", UpdatedAt: "2024-01-02 11:00"},
	{ID: 3, Title: "Refactoring", Author: "Martin Fowler", ISBN: "978-0134757599", Stock: 3, CreatedAt: "2024-01-03 12:00", UpdatedAt: "2024-01-03 12:00"},
}

// ListBooks 書籍清單頁
// c.HTML() 三個參數：HTTP 狀態碼、模板名稱、傳入模板的資料
func ListBooks(c *gin.Context) {
	c.HTML(http.StatusOK, "book/list.html", gin.H{
		"title": "書籍清單",
		"books": fakeBooks,
	})
}

// ShowBook 書籍詳情頁
func ShowBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "無效的 ID")
		return
	}

	// 從假資料中找書（後續 U06 改為查 DB）
	var book *Book
	for _, b := range fakeBooks {
		if b.ID == id {
			book = &b
			break
		}
	}

	if book == nil {
		c.String(http.StatusNotFound, "找不到書籍")
		return
	}

	c.HTML(http.StatusOK, "book/detail.html", gin.H{
		"title": book.Title,
		"book":  book,
	})
}

// NewBookForm 顯示新增書籍表單
func NewBookForm(c *gin.Context) {
	c.HTML(http.StatusOK, "book/form.html", gin.H{
		"title":  "新增書籍",
		"isEdit": false,
		"book":   nil,
		"errors": nil,
	})
}

// CreateBook 接收新增書籍表單（目前只是假實作，後續 U06 會存入 DB）
func CreateBook(c *gin.Context) {
	// 目前只是重導回清單頁，後續 U06 會實作真正的新增邏輯
	c.Redirect(http.StatusFound, "/books")
}

// EditBookForm 顯示編輯書籍表單
func EditBookForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "無效的 ID")
		return
	}

	var book *Book
	for _, b := range fakeBooks {
		if b.ID == id {
			book = &b
			break
		}
	}

	if book == nil {
		c.String(http.StatusNotFound, "找不到書籍")
		return
	}

	c.HTML(http.StatusOK, "book/form.html", gin.H{
		"title":  "編輯書籍",
		"isEdit": true,
		"book":   book,
		"errors": nil,
	})
}

// UpdateBook 接收編輯書籍表單（目前只是假實作，後續 U06 會更新 DB）
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	c.Redirect(http.StatusFound, "/books/"+id)
}

// DeleteBook 刪除書籍（目前只是假實作，後續 U06 會從 DB 刪除）
func DeleteBook(c *gin.Context) {
	c.Redirect(http.StatusFound, "/books")
}
