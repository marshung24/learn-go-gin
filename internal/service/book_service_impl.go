package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/example/learn-go-gin/internal/cache"
	"github.com/example/learn-go-gin/internal/model"
	"github.com/example/learn-go-gin/internal/model/dto"
	"github.com/example/learn-go-gin/internal/repository"
	"github.com/redis/go-redis/v9"
)

// bookServiceImpl 是 BookService 的 Cache-Aside 實作
// 讀取流程：先查 Redis，miss 時查 DB 並回寫快取（TTL 10 分鐘）
// 寫入/刪除流程：操作 DB 後驅逐快取，確保下次讀取重查 DB
type bookServiceImpl struct {
	repo        repository.BookRepository
	redisClient *redis.Client
}

// NewBookService 建立 BookService 實例
func NewBookService(repo repository.BookRepository, redisClient *redis.Client) BookService {
	return &bookServiceImpl{
		repo:        repo,
		redisClient: redisClient,
	}
}

// FindAll 查詢全部書籍（不走快取，避免全表快取的複雜性）
func (s *bookServiceImpl) FindAll() ([]dto.BookResponseDTO, error) {
	books, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.FromBooks(books), nil
}

// FindByID 依 ID 查詢書籍（Cache-Aside 模式）
func (s *bookServiceImpl) FindByID(id uint) (*dto.BookResponseDTO, error) {
	key := fmt.Sprintf("%s%d", cache.KeyPrefix, id)

	// Step 1: 查 Redis 快取
	data, err := s.redisClient.Get(cache.Ctx, key).Bytes()
	if err == nil {
		// Cache hit
		var response dto.BookResponseDTO
		if err := json.Unmarshal(data, &response); err == nil {
			log.Printf("Cache hit: %s", key)
			return &response, nil
		}
	}

	// Step 2: Cache miss → 查 DB
	log.Printf("Cache miss: %s", key)
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Step 3: 回寫快取
	response := dto.FromBook(book)
	if jsonData, err := json.Marshal(response); err == nil {
		s.redisClient.Set(cache.Ctx, key, jsonData, cache.CacheTTL)
	}

	return response, nil
}

// FindByISBN 依 ISBN 查詢書籍（不走快取，避免複雜的 key 管理）
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

// Update 更新書籍資料並驅逐快取
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

	// 驅逐快取，確保下次讀取重查 DB
	key := fmt.Sprintf("%s%d", cache.KeyPrefix, id)
	s.redisClient.Del(cache.Ctx, key)
	log.Printf("Cache evicted: %s", key)

	// 重新查詢以取得更新後的時間戳
	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return dto.FromBook(updated), nil
}

// Delete 刪除書籍並驅逐快取
func (s *bookServiceImpl) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 驅逐快取
	key := fmt.Sprintf("%s%d", cache.KeyPrefix, id)
	s.redisClient.Del(cache.Ctx, key)
	log.Printf("Cache evicted: %s", key)

	return nil
}
