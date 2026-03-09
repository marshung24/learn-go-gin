# U06｜GORM 串接資料庫

> 能用 GORM 完成真實 CRUD，頁面與 API 顯示 DB 真實資料 ｜ 120 min ｜ 前置依賴：U02, U03, U05

---

## 為什麼先教這個？

U05 設計好了資料表，這堂課要把假資料換成真的——用 GORM 串接 MySQL，讓 Handler 直接從 DB 讀寫資料。**完成後學員會第一次看到「頁面上的資料來自真實資料庫」的完整 CRUD 循環。** 此刻先不抽 Service 層，Handler 直接操作 Repository——保持最短路徑讓 CRUD 跑起來。

> **講師提示**：這是課程的第一個「啊哈時刻」——學員親眼看到頁面從假資料切換到真實 DB 資料。建議 live demo 時先秀假資料頁面，再切換到 Repository 呼叫，重新整理瀏覽器的瞬間最有衝擊力。

---

## 對應檔案

| 檔案 | 角色 |
|------|------|
| `model/book.go` | DB Model（對應 books 表） |
| `repository/book_repository.go` | Repository 介面定義 |
| `repository/book_repository_impl.go` | GORM 實作 |

---

## 核心觀念

### GORM Model 對應

```go
type Book struct {
    ID        uint           `gorm:"primaryKey"`
    Title     string         `gorm:"column:title;size:200;not null"`
    Author    string         `gorm:"column:author;size:100;not null"`
    ISBN      string         `gorm:"column:isbn;size:20;uniqueIndex"`
    Stock     int            `gorm:"column:stock;default:0"`
    CreatedAt time.Time      `gorm:"column:created_at"`
    UpdatedAt time.Time      `gorm:"column:updated_at"`
}
```

- `gorm:"primaryKey"` — 標示主鍵
- `gorm:"column:xxx"` — 對應 DB 欄位名稱
- `gorm:"uniqueIndex"` — 建立唯一索引

### GORM 自動遷移

```go
db.AutoMigrate(&Book{})
```

可自動同步 schema（教學用），**生產環境建議用 golang-migrate**。

### Repository Pattern

定義介面，實作與介面分離，方便測試時 mock：

```go
// 介面定義
type BookRepository interface {
    FindAll() ([]Book, error)
    FindByID(id uint) (*Book, error)
    Create(book *Book) error
    Update(book *Book) error
    Delete(id uint) error
}

// 實作
type bookRepositoryImpl struct {
    db *gorm.DB
}

func (r *bookRepositoryImpl) FindAll() ([]Book, error) {
    var books []Book
    result := r.db.Find(&books)
    return books, result.Error
}
```

### `db.First()` vs `db.Find()`

```go
// First 取一筆（找不到會 error）
var book Book
db.First(&book, id)  // 找不到會回傳 gorm.ErrRecordNotFound

// Find 取多筆（找不到回傳空 slice）
var books []Book
db.Find(&books)  // 空結果不會 error
```

### 替換假資料

修改 U02/U03 的 Handler，把 hardcoded slice 改成呼叫 repository：

```go
// 修改前（假資料）
books := []map[string]any{
    {"ID": 1, "Title": "Clean Code", ...},
}

// 修改後（真實 DB）
books, err := repo.FindAll()
if err != nil {
    // handle error
}
```

頁面立刻顯示真實 DB 資料。

---

## 動手做

### 必做

**1. 閱讀 Repository 結構**

閱讀 `book_repository.go` 介面和 `book_repository_impl.go` 實作，理解每個方法對應的 GORM 操作。

**2. 替換假資料**

修改 `book_view.go`，把 `ListBooks()` 方法中的假資料替換為 `repo.FindAll()`：

```go
func (h *BookViewHandler) ListBooks(c *gin.Context) {
    books, err := h.repo.FindAll()
    if err != nil {
        c.String(http.StatusInternalServerError, "Error fetching books")
        return
    }
    c.HTML(http.StatusOK, "book/list.html", gin.H{
        "books": books,
    })
}
```

重新整理頁面確認顯示 DB 資料。

**3. 新增 Repository 方法**

在 Repository 新增 `FindByISBN(isbn string) (*Book, error)` 方法：

```go
func (r *bookRepositoryImpl) FindByISBN(isbn string) (*Book, error) {
    var book Book
    result := r.db.Where("isbn = ?", isbn).First(&book)
    if result.Error != nil {
        return nil, result.Error
    }
    return &book, nil
}
```

確認編譯通過。

### 延伸挑戰

1. 新增一個多條件查詢方法（如依作者 + 庫存範圍），體驗 GORM 的 `Where` 鏈式呼叫：

```go
func (r *bookRepositoryImpl) FindByAuthorAndStock(author string, minStock int) ([]Book, error) {
    var books []Book
    result := r.db.Where("author = ?", author).Where("stock >= ?", minStock).Find(&books)
    return books, result.Error
}
```

---

## 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `record not found` panic | 用 `First()` 查不存在的資料時會 error | 檢查 `errors.Is(err, gorm.ErrRecordNotFound)` 或改用 `Find()` |
| 欄位值全是零值 | Model struct 的欄位名稱與 DB column 不對應 | 加上正確的 `gorm:"column:xxx"` tag |
| 連線失敗 | DSN 格式錯誤或 DB 服務未啟動 | 確認 docker compose 服務狀態，檢查 DSN 格式 |

---

## 驗收標準（DoD）

- [ ] 能在瀏覽器看到來自 DB 的書籍清單（不再是假資料）
- [ ] 能新增一個 Repository 方法
- [ ] 能解釋 `db.First()` 與 `db.Find()` 的差異
