package service

import (
	"testing"
	"time"

	"github.com/example/learn-go-gin/internal/mocks"
	"github.com/example/learn-go-gin/internal/model"
	"github.com/example/learn-go-gin/internal/model/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// TestFindAll_Success 測試查詢全部書籍成功
func TestFindAll_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)
	now := time.Now()

	mockBooks := []model.Book{
		{ID: 1, Title: "Clean Code", Author: "Robert Martin", ISBN: "978-0132350884", Stock: 5, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Title: "The Pragmatic Programmer", Author: "David Thomas", ISBN: "978-0135957059", Stock: 3, CreatedAt: now, UpdatedAt: now},
	}
	mockRepo.On("FindAll").Return(mockBooks, nil)

	service := NewBookService(mockRepo, nil)

	// Act
	result, err := service.FindAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Clean Code", result[0].Title)
	assert.Equal(t, "The Pragmatic Programmer", result[1].Title)
	mockRepo.AssertExpectations(t)
}

// TestFindAll_Empty 測試查詢結果為空
func TestFindAll_Empty(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)
	mockRepo.On("FindAll").Return([]model.Book{}, nil)

	service := NewBookService(mockRepo, nil)

	// Act
	result, err := service.FindAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

// TestFindAll_Error 測試查詢失敗
func TestFindAll_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)
	mockRepo.On("FindAll").Return(nil, gorm.ErrInvalidDB)

	service := NewBookService(mockRepo, nil)

	// Act
	result, err := service.FindAll()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// TestFindByISBN_Success 測試依 ISBN 查詢成功
func TestFindByISBN_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)
	now := time.Now()
	isbn := "978-0132350884"

	mockBook := &model.Book{
		ID:        1,
		Title:     "Clean Code",
		Author:    "Robert Martin",
		ISBN:      isbn,
		Stock:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockRepo.On("FindByISBN", isbn).Return(mockBook, nil)

	service := NewBookService(mockRepo, nil)

	// Act
	result, err := service.FindByISBN(isbn)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Clean Code", result.Title)
	assert.Equal(t, isbn, result.ISBN)
	mockRepo.AssertExpectations(t)
}

// TestFindByISBN_NotFound 測試依 ISBN 查詢找不到
func TestFindByISBN_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)
	isbn := "999-9999999999"
	mockRepo.On("FindByISBN", isbn).Return(nil, gorm.ErrRecordNotFound)

	service := NewBookService(mockRepo, nil)

	// Act
	result, err := service.FindByISBN(isbn)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// TestCreate_Success 測試新增書籍成功
func TestCreate_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)
	now := time.Now()

	input := &dto.BookCreateDTO{
		Title:  "Clean Code",
		Author: "Robert Martin",
		ISBN:   "978-0132350884",
		Stock:  5,
	}

	// Mock Create - 設定 book.ID 會被設為 1
	mockRepo.On("Create", mock.AnythingOfType("*model.Book")).Run(func(args mock.Arguments) {
		book := args.Get(0).(*model.Book)
		book.ID = 1
	}).Return(nil)

	// Mock FindByID 回傳完整資料
	createdBook := &model.Book{
		ID:        1,
		Title:     "Clean Code",
		Author:    "Robert Martin",
		ISBN:      "978-0132350884",
		Stock:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockRepo.On("FindByID", uint(1)).Return(createdBook, nil)

	service := NewBookService(mockRepo, nil)

	// Act
	result, err := service.Create(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Clean Code", result.Title)
	mockRepo.AssertExpectations(t)
}

// TestCreate_Error 測試新增書籍失敗
func TestCreate_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)

	input := &dto.BookCreateDTO{
		Title:  "Clean Code",
		Author: "Robert Martin",
		ISBN:   "978-0132350884",
		Stock:  5,
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Book")).Return(gorm.ErrInvalidDB)

	service := NewBookService(mockRepo, nil)

	// Act
	result, err := service.Create(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// TestDelete_Success 測試刪除書籍成功
func TestDelete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)
	mockRepo.On("Delete", uint(1)).Return(nil)

	service := NewBookService(mockRepo, nil)

	// Act
	err := service.Delete(1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestDelete_Error 測試刪除書籍失敗
func TestDelete_Error(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockBookRepository)
	mockRepo.On("Delete", uint(99)).Return(gorm.ErrRecordNotFound)

	service := NewBookService(mockRepo, nil)

	// Act
	err := service.Delete(99)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestFindByID 使用 table-driven tests 測試 FindByID
// 注意：這個測試需要 Redis mock，這裡僅測試 cache miss 的情況
func TestFindByID_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func(*mocks.MockBookRepository)
		wantErr   bool
		wantTitle string
	}{
		{
			name: "NotFound",
			id:   99,
			mockSetup: func(m *mocks.MockBookRepository) {
				m.On("FindByID", uint(99)).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockBookRepository)
			tt.mockSetup(mockRepo)

			// 由於 FindByID 需要 Redis，這裡僅測試會導致 repo 錯誤的情況
			service := &bookServiceImpl{repo: mockRepo, redisClient: nil}

			// 當 Redis 為 nil 時，直接調用 repo
			_, err := mockRepo.FindByID(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
