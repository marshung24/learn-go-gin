package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/example/learn-go-gin/internal/model/dto"
	"github.com/example/learn-go-gin/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookAPIHandler 書籍 REST API Handler
type BookAPIHandler struct {
	service service.BookService
}

// NewBookAPIHandler 建立 BookAPIHandler 實例
func NewBookAPIHandler(svc service.BookService) *BookAPIHandler {
	return &BookAPIHandler{service: svc}
}

// GetBooks godoc
// @Summary     取得所有書籍
// @Description 取得書籍清單
// @Tags        books
// @Produce     json
// @Success     200 {array} dto.BookResponseDTO
// @Router      /api/books [get]
func (h *BookAPIHandler) GetBooks(c *gin.Context) {
	books, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}

	// 確保回傳空陣列而非 null
	if books == nil {
		books = []dto.BookResponseDTO{}
	}

	c.JSON(http.StatusOK, books)
}

// GetBookByID godoc
// @Summary     依 ID 取得書籍
// @Description 依書籍 ID 取得單筆書籍資料
// @Tags        books
// @Produce     json
// @Param       id path int true "書籍 ID"
// @Success     200 {object} dto.BookResponseDTO
// @Failure     404 {object} map[string]string
// @Router      /api/books/{id} [get]
func (h *BookAPIHandler) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 ID"})
		return
	}

	book, err := h.service.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到書籍"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary     新增書籍
// @Description 新增一本書籍
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       book body dto.BookCreateDTO true "書籍資料"
// @Success     201 {object} dto.BookResponseDTO
// @Failure     400 {object} map[string]string
// @Router      /api/books [post]
func (h *BookAPIHandler) CreateBook(c *gin.Context) {
	var input dto.BookCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.service.Create(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// UpdateBook godoc
// @Summary     更新書籍
// @Description 依 ID 更新書籍資料
// @Tags        books
// @Accept      json
// @Produce     json
// @Param       id path int true "書籍 ID"
// @Param       book body dto.BookCreateDTO true "書籍資料"
// @Success     200 {object} dto.BookResponseDTO
// @Failure     404 {object} map[string]string
// @Router      /api/books/{id} [put]
func (h *BookAPIHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 ID"})
		return
	}

	var input dto.BookCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.service.Update(uint(id), &input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到書籍"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
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
	_, err = h.service.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到書籍"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting book"})
		return
	}

	c.Status(http.StatusNoContent)
}
