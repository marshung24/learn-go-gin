# U09｜Redis 快取整合

> 能用 redis-cli 驗證快取行為並說明 Cache-Aside 流程 ｜ 90 min ｜ 前置依賴：U07

---

## 為什麼先教這個？

資料庫查詢是 Web 應用最慢的環節。把「熱資料」放在 Redis 記憶體中，回應速度可以從毫秒降到微秒等級。這堂課在 U07 建立的 Service 層上加入快取邏輯。

---

## 對應檔案

| 檔案 | 角色 |
|------|------|
| `service/book_service_impl.go` | Cache-Aside 讀寫邏輯的實際所在 |
| `cache/redis.go` | Redis client 初始化與設定 |
| `cache/warmup.go` | 啟動時快取預熱（preload 常用資料） |

---

## 核心觀念

### Cache-Aside 模式

**讀取流程**：
```
1. 查 Redis 快取
2. 若 hit → 直接回傳
3. 若 miss → 查 DB → 寫入 Redis → 回傳
```

**寫入/刪除流程**：
```
1. 操作 DB（INSERT / UPDATE / DELETE）
2. 刪除對應的 Redis 快取 key
```

### go-redis 基本操作

```go
import "github.com/redis/go-redis/v9"

// 初始化
rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// 讀取
val, err := rdb.Get(ctx, "book:1").Result()
if err == redis.Nil {
    // key 不存在（cache miss）
}

// 寫入（含 TTL）
rdb.Set(ctx, "book:1", jsonData, 10*time.Minute)

// 刪除
rdb.Del(ctx, "book:1")
```

### JSON 序列化

```go
// Model → JSON string → Redis
data, _ := json.Marshal(book)
rdb.Set(ctx, key, data, ttl)

// Redis → JSON string → Model
data, _ := rdb.Get(ctx, key).Bytes()
var book Book
json.Unmarshal(data, &book)
```

### 快取 key 設計

```
book:{id}
```

例如：`book:1`、`book:42`

TTL（Time To Live）設為 10 分鐘。

### 為什麼 update/delete 要清快取？

避免讀到過期資料（快取一致性）：

```go
func (s *bookServiceImpl) Update(id uint, dto *BookCreateDTO) (*BookResponseDTO, error) {
    // 1. 更新 DB
    book, err := s.repo.Update(id, dto.ToModel())
    if err != nil {
        return nil, err
    }

    // 2. 刪除快取（確保下次讀取會從 DB 拿最新資料）
    s.cache.Del(ctx, fmt.Sprintf("book:%d", id))

    return book.ToResponse(), nil
}
```

### 預熱（Warmup）

應用啟動時預先載入常用資料到快取，減少 cold start 的 cache miss：

```go
func Warmup(repo repository.BookRepository, cache *redis.Client) {
    books, _ := repo.FindAll()
    for _, book := range books {
        data, _ := json.Marshal(book)
        cache.Set(ctx, fmt.Sprintf("book:%d", book.ID), data, 10*time.Minute)
    }
    log.Println("Cache warmup completed")
}
```

---

## 動手做

### 必做

**1. 追蹤 FindByID 快取流程**

追蹤 `book_service_impl.go` 的 `FindByID()` 完整流程：

```go
func (s *bookServiceImpl) GetByID(id uint) (*BookResponseDTO, error) {
    key := fmt.Sprintf("book:%d", id)

    // 1. 查 Redis
    data, err := s.cache.Get(ctx, key).Bytes()
    if err == nil {
        // cache hit
        var book Book
        json.Unmarshal(data, &book)
        return book.ToResponse(), nil
    }

    // 2. cache miss → 查 DB
    book, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // 3. 寫回 Redis
    data, _ = json.Marshal(book)
    s.cache.Set(ctx, key, data, 10*time.Minute)

    return book.ToResponse(), nil
}
```

**2. 用 redis-cli 觀察快取**

```bash
# 連線到 Redis
redis-cli

# 查看所有 book 相關的 key（教學環境限定）
KEYS book:*

# 查看某個 key 的值
GET book:1

# 查看 TTL
TTL book:1
```

**3. 驗證 cache miss 與回寫**

```bash
# 手動刪除一個 key
DEL book:1

# 呼叫 API
curl http://localhost:8080/api/books/1

# 觀察 console log 中的 cache miss 訊息

# 再次查看 Redis
GET book:1  # 應該有值了
```

### 延伸挑戰

1. 呼叫 `GET /api/books/1` 兩次，比較回應時間差異
2. 將 TTL 改為 1 分鐘觀察快取命中率變化
3. 思考：什麼情況下應該用 Cache-Aside？什麼情況下不適合？

---

## 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `KEYS` 指令在正式環境不能用 | `KEYS` 會阻塞 Redis（全表掃描） | 教學環境直觀好用；正式環境改用 `SCAN` |
| 快取與 DB 資料不一致 | 更新 DB 後忘了清快取 | 確保 Update / Delete 操作都有對應的 `rdb.Del(ctx, key)` |
| Redis 連線失敗 | 容器未啟動或設定錯誤 | 檢查 docker compose 狀態，確認 Redis host/port 設定 |

---

## 驗收標準（DoD）

- [ ] 能用 `redis-cli` 驗證快取行為（確認 key 存在、手動刪除後觀察回源）
- [ ] 能說明 Cache-Aside 的讀寫流程
- [ ] 能說明為什麼 Update/Delete 後要清快取
