package dto

import (
	"time"

	"github.com/example/learn-go-gin/internal/model"
)

// BookResponseDTO 書籍查詢回應 DTO
// 包含完整欄位（含 DB 產生的時間戳）
type BookResponseDTO struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// FromBook 從 Book Model 轉換為 ResponseDTO
func FromBook(book *model.Book) *BookResponseDTO {
	return &BookResponseDTO{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		ISBN:      book.ISBN,
		Stock:     book.Stock,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
}

// FromBooks 從 Book slice 轉換為 ResponseDTO slice
func FromBooks(books []model.Book) []BookResponseDTO {
	result := make([]BookResponseDTO, len(books))
	for i, book := range books {
		result[i] = *FromBook(&book)
	}
	return result
}
