package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/example/learn-go-gin/internal/model/dto"
	"github.com/example/learn-go-gin/internal/repository"
)

// KeyPrefix 快取 key 前綴
const KeyPrefix = "book:"

// CacheTTL 快取存活時間（10 分鐘）
const CacheTTL = 10 * time.Minute

// Warmup 快取預熱
// 應用啟動時預先載入全部書籍到 Redis，減少 cold start 的 cache miss
func Warmup(repo repository.BookRepository) {
	books, err := repo.FindAll()
	if err != nil {
		log.Printf("Cache warmup failed: %v", err)
		return
	}

	for _, book := range books {
		response := dto.FromBook(&book)
		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal book %d: %v", book.ID, err)
			continue
		}

		key := fmt.Sprintf("%s%d", KeyPrefix, book.ID)
		if err := RedisClient.Set(Ctx, key, data, CacheTTL).Err(); err != nil {
			log.Printf("Failed to cache book %d: %v", book.ID, err)
			continue
		}
	}

	log.Printf("Cache warmup complete. %d books loaded.", len(books))
}
