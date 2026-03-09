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

// BookAPIHandler 書籍 REST API Handler
type BookAPIHandler struct {
	repo repository.BookRepository
}

// NewBookAPIHandler 建立 BookAPIHandler 實例
func NewBookAPIHandler(repo repository.BookRepository) *BookAPIHandler {
	return &BookAPIHandler{repo: repo}
}

// BookResponse 是 API 回應的書籍資料結構
type BookResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	ISBN      string `json:"isbn"`
	Stock     int    `json:"stock"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// BookCreateRequest 是新增/更新書籍的請求結構
type BookCreateRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	ISBN   string `json:"isbn" binding:"required"`
	Stock  int    `json:"stock" binding:"min=0"`
}

// toResponse 將 model.Book 轉換為 API 回應格式
func toResponse(book *model.Book) BookResponse {
	return BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		ISBN:      book.ISBN,
		Stock:     book.Stock,
		CreatedAt: book.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: book.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// GetBooks godoc
// @Summary     取得所有書籍
// @Description 取得書籍清單
// @Tags        books
// @Produce     json
// @Success     200 {array} BookResponse
// @Router      /api/books [get]
func (h *BookAPIHandler) GetBooks(c *gin.Context) {
	books, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}

	var responses []BookResponse
	for _, book := range books {
		responses = append(responses, toResponse(&book))
	}

	// 確保回傳空陣列而非 null
	if responses == nil {
		responses = []BookResponse{}
	}

	c.JSON(http.StatusOK, responses)
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
func (h *BookAPIHandler) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 ID"})
		return
	}

	book, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到書籍"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		return
	}

	c.JSON(http.StatusOK, toResponse(book))
}

// CreateBook godoc
// @Summary     新增書籍
// @Description 新增一本書籍
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       book body BookCreateRequest true "書籍資料"
// @Success     201 {object} BookResponse
// @Failure     400 {object} map[string]string
// @Router      /api/books [post]
func (h *BookAPIHandler) CreateBook(c *gin.Context) {
	var req BookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := &model.Book{
		Title:  req.Title,
		Author: req.Author,
		ISBN:   req.ISBN,
		Stock:  req.Stock,
	}

	if err := h.repo.Create(book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 重新查詢以取得完整資料（包含時間戳）
	created, _ := h.repo.FindByID(book.ID)
	c.JSON(http.StatusCreated, toResponse(created))
}

// UpdateBook godoc
// @Summary     更新書籍
// @Description 依 ID 更新書籍資料
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       id path int true "書籍 ID"
// @Param       book body BookCreateRequest true "書籍資料"
// @Success     200 {object} BookResponse
// @Failure     404 {object} map[string]string
// @Router      /api/books/{id} [put]
func (h *BookAPIHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 ID"})
		return
	}

	book, err := h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到書籍"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		return
	}

	var req BookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.Title = req.Title
	book.Author = req.Author
	book.ISBN = req.ISBN
	book.Stock = req.Stock

	if err := h.repo.Update(book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 重新查詢以取得更新後的時間戳
	updated, _ := h.repo.FindByID(book.ID)
	c.JSON(http.StatusOK, toResponse(updated))
}

// DeleteBook godoc
// @Summary     刪除書籍
// @Description 依 ID 刪除書籍
// @Tags        books
// @Param       id path int true "書籍 ID"
// @Success     204
// @Failure     404 {object} map[string]string
// @Router      /api/books/{id} [delete]
func (h *BookAPIHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 ID"})
		return
	}

	// 先確認存在
	_, err = h.repo.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到書籍"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting book"})
		return
	}

	c.Status(http.StatusNoContent)
}
