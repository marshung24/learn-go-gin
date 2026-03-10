package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/learn-go-gin/internal/mocks"
	"github.com/example/learn-go-gin/internal/model/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func init() {
	// 設定 Gin 為測試模式，減少日誌輸出
	gin.SetMode(gin.TestMode)
}

// setupRouter 建立測試用的 Gin router
func setupRouter(handler *BookAPIHandler) *gin.Engine {
	r := gin.New()
	r.GET("/api/books", handler.GetBooks)
	r.GET("/api/books/:id", handler.GetBookByID)
	r.POST("/api/books", handler.CreateBook)
	r.PUT("/api/books/:id", handler.UpdateBook)
	r.DELETE("/api/books/:id", handler.DeleteBook)
	return r
}

// TestGetBooks_Success 測試取得所有書籍成功
func TestGetBooks_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	now := time.Now()

	mockBooks := []dto.BookResponseDTO{
		{ID: 1, Title: "Clean Code", Author: "Robert Martin", ISBN: "978-0132350884", Stock: 5, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Title: "The Pragmatic Programmer", Author: "David Thomas", ISBN: "978-0135957059", Stock: 3, CreatedAt: now, UpdatedAt: now},
	}
	mockService.On("FindAll").Return(mockBooks, nil)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var books []dto.BookResponseDTO
	err := json.Unmarshal(rec.Body.Bytes(), &books)
	assert.NoError(t, err)
	assert.Len(t, books, 2)
	assert.Equal(t, "Clean Code", books[0].Title)

	mockService.AssertExpectations(t)
}

// TestGetBooks_Empty 測試取得空書籍清單
func TestGetBooks_Empty(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	mockService.On("FindAll").Return([]dto.BookResponseDTO{}, nil)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var books []dto.BookResponseDTO
	err := json.Unmarshal(rec.Body.Bytes(), &books)
	assert.NoError(t, err)
	assert.Len(t, books, 0)

	mockService.AssertExpectations(t)
}

// TestGetBooks_Error 測試取得書籍時發生錯誤
func TestGetBooks_Error(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	mockService.On("FindAll").Return(nil, gorm.ErrInvalidDB)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Error fetching books")

	mockService.AssertExpectations(t)
}

// TestGetBookByID_Success 測試依 ID 取得書籍成功
func TestGetBookByID_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	now := time.Now()

	mockBook := &dto.BookResponseDTO{
		ID:        1,
		Title:     "Clean Code",
		Author:    "Robert Martin",
		ISBN:      "978-0132350884",
		Stock:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockService.On("FindByID", uint(1)).Return(mockBook, nil)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/books/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var book dto.BookResponseDTO
	err := json.Unmarshal(rec.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), book.ID)
	assert.Equal(t, "Clean Code", book.Title)

	mockService.AssertExpectations(t)
}

// TestGetBookByID_NotFound 測試依 ID 取得書籍找不到
func TestGetBookByID_NotFound(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	mockService.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/books/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "找不到書籍")

	mockService.AssertExpectations(t)
}

// TestGetBookByID_InvalidID 測試無效的 ID
func TestGetBookByID_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/books/abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "無效的 ID")
}

// TestCreateBook_Success 測試新增書籍成功
func TestCreateBook_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	now := time.Now()

	input := dto.BookCreateDTO{
		Title:  "Clean Code",
		Author: "Robert Martin",
		ISBN:   "978-0132350884",
		Stock:  5,
	}

	createdBook := &dto.BookResponseDTO{
		ID:        1,
		Title:     "Clean Code",
		Author:    "Robert Martin",
		ISBN:      "978-0132350884",
		Stock:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockService.On("Create", mock.AnythingOfType("*dto.BookCreateDTO")).Return(createdBook, nil)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

	var book dto.BookResponseDTO
	err := json.Unmarshal(rec.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), book.ID)
	assert.Equal(t, "Clean Code", book.Title)

	mockService.AssertExpectations(t)
}

// TestCreateBook_ValidationError 測試新增書籍驗證失敗
func TestCreateBook_ValidationError(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// 缺少必填欄位
	input := map[string]interface{}{
		"title": "", // 空白標題
		"stock": 5,
	}

	// Act
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["error"])
}

// TestUpdateBook_Success 測試更新書籍成功
func TestUpdateBook_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	now := time.Now()

	input := dto.BookCreateDTO{
		Title:  "Clean Code (2nd Edition)",
		Author: "Robert Martin",
		ISBN:   "978-0132350884",
		Stock:  10,
	}

	updatedBook := &dto.BookResponseDTO{
		ID:        1,
		Title:     "Clean Code (2nd Edition)",
		Author:    "Robert Martin",
		ISBN:      "978-0132350884",
		Stock:     10,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.BookCreateDTO")).Return(updatedBook, nil)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPut, "/api/books/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var book dto.BookResponseDTO
	err := json.Unmarshal(rec.Body.Bytes(), &book)
	assert.NoError(t, err)
	assert.Equal(t, "Clean Code (2nd Edition)", book.Title)
	assert.Equal(t, 10, book.Stock)

	mockService.AssertExpectations(t)
}

// TestUpdateBook_NotFound 測試更新不存在的書籍
func TestUpdateBook_NotFound(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)

	input := dto.BookCreateDTO{
		Title:  "Clean Code",
		Author: "Robert Martin",
		ISBN:   "978-0132350884",
		Stock:  5,
	}

	mockService.On("Update", uint(99), mock.AnythingOfType("*dto.BookCreateDTO")).Return(nil, gorm.ErrRecordNotFound)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPut, "/api/books/99", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "找不到書籍")

	mockService.AssertExpectations(t)
}

// TestDeleteBook_Success 測試刪除書籍成功
func TestDeleteBook_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	now := time.Now()

	mockBook := &dto.BookResponseDTO{
		ID:        1,
		Title:     "Clean Code",
		Author:    "Robert Martin",
		ISBN:      "978-0132350884",
		Stock:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockService.On("FindByID", uint(1)).Return(mockBook, nil)
	mockService.On("Delete", uint(1)).Return(nil)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodDelete, "/api/books/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

// TestDeleteBook_NotFound 測試刪除不存在的書籍
func TestDeleteBook_NotFound(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	mockService.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)

	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodDelete, "/api/books/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "找不到書籍")

	mockService.AssertExpectations(t)
}

// TestDeleteBook_InvalidID 測試刪除無效 ID
func TestDeleteBook_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(mocks.MockBookService)
	handler := NewBookAPIHandler(mockService)
	router := setupRouter(handler)

	// Act
	req := httptest.NewRequest(http.MethodDelete, "/api/books/abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "無效的 ID")
}
