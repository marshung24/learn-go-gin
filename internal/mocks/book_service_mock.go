package mocks

import (
	"github.com/example/learn-go-gin/internal/model/dto"
	"github.com/stretchr/testify/mock"
)

// MockBookService 是 BookService 的 mock 實作
type MockBookService struct {
	mock.Mock
}

// FindAll mock 實作
func (m *MockBookService) FindAll() ([]dto.BookResponseDTO, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.BookResponseDTO), args.Error(1)
}

// FindByID mock 實作
func (m *MockBookService) FindByID(id uint) (*dto.BookResponseDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BookResponseDTO), args.Error(1)
}

// FindByISBN mock 實作
func (m *MockBookService) FindByISBN(isbn string) (*dto.BookResponseDTO, error) {
	args := m.Called(isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BookResponseDTO), args.Error(1)
}

// Create mock 實作
func (m *MockBookService) Create(input *dto.BookCreateDTO) (*dto.BookResponseDTO, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BookResponseDTO), args.Error(1)
}

// Update mock 實作
func (m *MockBookService) Update(id uint, input *dto.BookCreateDTO) (*dto.BookResponseDTO, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BookResponseDTO), args.Error(1)
}

// Delete mock 實作
func (m *MockBookService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
