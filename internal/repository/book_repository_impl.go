package repository

import (
	"github.com/example/learn-go-gin/internal/model"
	"gorm.io/gorm"
)

// bookRepositoryImpl 是 BookRepository 的 GORM 實作
type bookRepositoryImpl struct {
	db *gorm.DB
}

// NewBookRepository 建立 BookRepository 實例
func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepositoryImpl{db: db}
}

// FindAll 查詢全部書籍，依 id 降冪排列
// db.Find() 找不到資料回傳空 slice，不會 error
func (r *bookRepositoryImpl) FindAll() ([]model.Book, error) {
	var books []model.Book
	result := r.db.Order("id desc").Find(&books)
	return books, result.Error
}

// FindByID 依主鍵查詢單筆書籍
// db.First() 找不到會回傳 gorm.ErrRecordNotFound
func (r *bookRepositoryImpl) FindByID(id uint) (*model.Book, error) {
	var book model.Book
	result := r.db.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

// FindByISBN 依 ISBN 查詢書籍
func (r *bookRepositoryImpl) FindByISBN(isbn string) (*model.Book, error) {
	var book model.Book
	result := r.db.Where("isbn = ?", isbn).First(&book)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

// Create 新增書籍
// GORM 會自動回填 DB 產生的 ID 到 book.ID
func (r *bookRepositoryImpl) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

// Update 更新書籍資料
// 使用 Save() 會更新所有欄位（包含零值）
func (r *bookRepositoryImpl) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

// Delete 依主鍵刪除書籍
func (r *bookRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&model.Book{}, id).Error
}
