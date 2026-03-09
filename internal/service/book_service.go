package service

import (
	"github.com/example/learn-go-gin/internal/model/dto"
)

// BookService 書籍業務邏輯服務介面
// 定義書籍 CRUD 操作契約
type BookService interface {
	// FindAll 查詢全部書籍，依 id 降冪排列
	FindAll() ([]dto.BookResponseDTO, error)

	// FindByID 依 ID 查詢書籍
	FindByID(id uint) (*dto.BookResponseDTO, error)

	// FindByISBN 依 ISBN 查詢書籍
	FindByISBN(isbn string) (*dto.BookResponseDTO, error)

	// Create 新增書籍
	Create(input *dto.BookCreateDTO) (*dto.BookResponseDTO, error)

	// Update 更新書籍資料
	Update(id uint, input *dto.BookCreateDTO) (*dto.BookResponseDTO, error)

	// Delete 刪除書籍
	Delete(id uint) error
}
