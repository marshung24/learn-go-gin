package repository

import (
	"github.com/example/learn-go-gin/internal/model"
)

// BookRepository 書籍資料存取介面
// 定義介面讓實作與介面分離，方便測試時 mock
type BookRepository interface {
	// FindAll 查詢全部書籍，依 id 降冪排列
	FindAll() ([]model.Book, error)

	// FindByID 依主鍵查詢單筆書籍
	FindByID(id uint) (*model.Book, error)

	// FindByISBN 依 ISBN 查詢書籍
	FindByISBN(isbn string) (*model.Book, error)

	// Create 新增書籍
	Create(book *model.Book) error

	// Update 更新書籍資料
	Update(book *model.Book) error

	// Delete 依主鍵刪除書籍
	Delete(id uint) error
}
