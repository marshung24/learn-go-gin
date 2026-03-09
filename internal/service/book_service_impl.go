package service

import (
	"github.com/example/learn-go-gin/internal/model"
	"github.com/example/learn-go-gin/internal/model/dto"
	"github.com/example/learn-go-gin/internal/repository"
)

// bookServiceImpl 是 BookService 的實作
// 此刻先不含快取，U09 會加入 Redis Cache-Aside
type bookServiceImpl struct {
	repo repository.BookRepository
}

// NewBookService 建立 BookService 實例
func NewBookService(repo repository.BookRepository) BookService {
	return &bookServiceImpl{repo: repo}
}

// FindAll 查詢全部書籍
func (s *bookServiceImpl) FindAll() ([]dto.BookResponseDTO, error) {
	books, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.FromBooks(books), nil
}

// FindByID 依 ID 查詢書籍
func (s *bookServiceImpl) FindByID(id uint) (*dto.BookResponseDTO, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return dto.FromBook(book), nil
}

// FindByISBN 依 ISBN 查詢書籍
func (s *bookServiceImpl) FindByISBN(isbn string) (*dto.BookResponseDTO, error) {
	book, err := s.repo.FindByISBN(isbn)
	if err != nil {
		return nil, err
	}
	return dto.FromBook(book), nil
}

// Create 新增書籍
func (s *bookServiceImpl) Create(input *dto.BookCreateDTO) (*dto.BookResponseDTO, error) {
	book := &model.Book{
		Title:  input.Title,
		Author: input.Author,
		ISBN:   input.ISBN,
		Stock:  input.Stock,
	}

	if err := s.repo.Create(book); err != nil {
		return nil, err
	}

	// 重新查詢以取得完整資料（包含時間戳）
	created, err := s.repo.FindByID(book.ID)
	if err != nil {
		return nil, err
	}

	return dto.FromBook(created), nil
}

// Update 更新書籍資料
func (s *bookServiceImpl) Update(id uint, input *dto.BookCreateDTO) (*dto.BookResponseDTO, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	book.Title = input.Title
	book.Author = input.Author
	book.ISBN = input.ISBN
	book.Stock = input.Stock

	if err := s.repo.Update(book); err != nil {
		return nil, err
	}

	// 重新查詢以取得更新後的時間戳
	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return dto.FromBook(updated), nil
}

// Delete 刪除書籍
func (s *bookServiceImpl) Delete(id uint) error {
	return s.repo.Delete(id)
}
