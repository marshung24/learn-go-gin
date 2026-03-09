package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/example/learn-go-gin/internal/model"
	"github.com/example/learn-go-gin/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookViewHandler 書籍 MVC Handler（頁面渲染）
type BookViewHandler struct {
	repo repository.BookRepository
}

// NewBookViewHandler 建立 BookViewHandler 實例
func NewBookViewHandler(repo repository.BookRepository) *BookViewHandler {
	return &BookViewHandler{repo: repo}
}

// ListBooks 書籍清單頁
// 從 Repository 取得真實 DB 資料
func (h *BookViewHandler) ListBooks(c *gin.Context) {
	books, err := h.repo.FindAll()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching books")
		return
	}

	c.HTML(http.StatusOK, "book/list.html", gin.H{
		"title": "書籍清單",
		"books": books,
	})
}

// ShowBook 書籍詳情頁
func (h *BookViewHandler) ShowBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "無效的 ID")
		return
	}

	book, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, "找不到書籍")
			return
		}
		c.String(http.StatusInternalServerError, "Error fetching book")
		return
	}

	c.HTML(http.StatusOK, "book/detail.html", gin.H{
		"title": book.Title,
		"book":  book,
	})
}

// NewBookForm 顯示新增書籍表單
func (h *BookViewHandler) NewBookForm(c *gin.Context) {
	c.HTML(http.StatusOK, "book/form.html", gin.H{
		"title":  "新增書籍",
		"isEdit": false,
		"book":   nil,
		"errors": nil,
	})
}

// CreateBook 接收新增書籍表單
func (h *BookViewHandler) CreateBook(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	isbn := c.PostForm("isbn")
	stockStr := c.PostForm("stock")

	stock, _ := strconv.Atoi(stockStr)

	book := &model.Book{
		Title:  title,
		Author: author,
		ISBN:   isbn,
		Stock:  stock,
	}

	if err := h.repo.Create(book); err != nil {
		c.HTML(http.StatusBadRequest, "book/form.html", gin.H{
			"title":  "新增書籍",
			"isEdit": false,
			"book":   book,
			"errors": map[string]string{"general": err.Error()},
		})
		return
	}

	c.Redirect(http.StatusFound, "/books")
}

// EditBookForm 顯示編輯書籍表單
func (h *BookViewHandler) EditBookForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "無效的 ID")
		return
	}

	book, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, "找不到書籍")
			return
		}
		c.String(http.StatusInternalServerError, "Error fetching book")
		return
	}

	c.HTML(http.StatusOK, "book/form.html", gin.H{
		"title":  "編輯書籍",
		"isEdit": true,
		"book":   book,
		"errors": nil,
	})
}

// UpdateBook 接收編輯書籍表單
func (h *BookViewHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "無效的 ID")
		return
	}

	book, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, "找不到書籍")
			return
		}
		c.String(http.StatusInternalServerError, "Error fetching book")
		return
	}

	// 更新欄位
	book.Title = c.PostForm("title")
	book.Author = c.PostForm("author")
	book.ISBN = c.PostForm("isbn")
	stock, _ := strconv.Atoi(c.PostForm("stock"))
	book.Stock = stock

	if err := h.repo.Update(book); err != nil {
		c.HTML(http.StatusBadRequest, "book/form.html", gin.H{
			"title":  "編輯書籍",
			"isEdit": true,
			"book":   book,
			"errors": map[string]string{"general": err.Error()},
		})
		return
	}

	c.Redirect(http.StatusFound, "/books/"+idStr)
}

// DeleteBook 刪除書籍
func (h *BookViewHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "無效的 ID")
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.String(http.StatusInternalServerError, "Error deleting book")
		return
	}

	c.Redirect(http.StatusFound, "/books")
}
