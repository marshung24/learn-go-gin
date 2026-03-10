package mocks

import (
	"github.com/example/learn-go-gin/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockBookRepository 是 BookRepository 的 mock 實作
type MockBookRepository struct {
	mock.Mock
}

// FindAll mock 實作
func (m *MockBookRepository) FindAll() ([]model.Book, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Book), args.Error(1)
}

// FindByID mock 實作
func (m *MockBookRepository) FindByID(id uint) (*model.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Book), args.Error(1)
}

// FindByISBN mock 實作
func (m *MockBookRepository) FindByISBN(isbn string) (*model.Book, error) {
	args := m.Called(isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Book), args.Error(1)
}

// Create mock 實作
func (m *MockBookRepository) Create(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

// Update mock 實作
func (m *MockBookRepository) Update(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

// Delete mock 實作
func (m *MockBookRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
