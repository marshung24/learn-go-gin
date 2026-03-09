# U07｜正規化架構 — Service 與 DTO 分層

> 能從 Handler 抽出 Service 層，引入 DTO 做資料驗證與格式轉換 ｜ 90 min ｜ 前置依賴：U06

---

## 為什麼先教這個？

U06 完成了 CRUD，但 Handler 直接操作 Repository——所有邏輯擠在一起，改一個地方可能影響全部。這堂課把程式碼「整理乾淨」：抽出 Service 層集中業務邏輯，引入 DTO 分離輸入/輸出格式。**重構前後功能完全一樣，但程式碼更好維護、更好測試。**

> **講師提示**：先展示重構前（Handler 直接操作 Repository）和重構後（Handler → Service → Repository）的對比，讓學員感受到「程式碼可以一樣能跑，但結構差很多」。

---

## 對應檔案

| 檔案 | 角色 |
|------|------|
| `service/book_service.go` | Service 介面（定義業務合約） |
| `service/book_service_impl.go` | Service 實作（業務邏輯，此刻先不含快取） |
| `model/dto/book_create_dto.go` | 輸入驗證 DTO（新增/更新用） |
| `model/dto/book_response_dto.go` | 輸出格式 DTO（API 回傳用） |

---

## 核心觀念

### 介面與實作分離

Go 的 implicit interface 讓測試可以用 Mock 替換實作：

```go
// 介面定義
type BookService interface {
    GetAll() ([]BookResponseDTO, error)
    GetByID(id uint) (*BookResponseDTO, error)
    Create(dto *BookCreateDTO) (*BookResponseDTO, error)
    Update(id uint, dto *BookCreateDTO) (*BookResponseDTO, error)
    Delete(id uint) error
}

// 實作
type bookServiceImpl struct {
    repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
    return &bookServiceImpl{repo: repo}
}
```

### 業務責任邊界

| 層級 | 職責 | 範例 |
|------|------|------|
| **Handler** | HTTP 協議 | 回傳 201 狀態碼、解析 JSON body |
| **Service** | 業務規則 | 檢查重複 ISBN、計算庫存 |
| **Repository** | 資料存取 | 執行 SQL 查詢 |

兩者不應互相滲透——Service 不該知道 HTTP 狀態碼，Handler 不該知道 SQL。

### Model 與 DTO 分離

```go
// Book — 對應 DB
type Book struct {
    ID        uint
    Title     string
    Author    string
    ISBN      string
    Stock     int
    CreatedAt time.Time
    UpdatedAt time.Time
}

// BookCreateDTO — 收什麼（輸入驗證）
type BookCreateDTO struct {
    Title  string `json:"title" binding:"required,min=1,max=200"`
    Author string `json:"author" binding:"required,min=1,max=100"`
    ISBN   string `json:"isbn" binding:"omitempty,len=13"`
    Stock  int    `json:"stock" binding:"gte=0"`
}

// BookResponseDTO — 給什麼（輸出格式）
type BookResponseDTO struct {
    ID     uint   `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    ISBN   string `json:"isbn"`
    Stock  int    `json:"stock"`
}
```

修改其中一邊不會牽動另一邊。

### Gin binding 驗證

`binding:"required,min=1"` 在 DTO struct tag 上宣告驗證規則：

```go
type BookCreateDTO struct {
    Title  string `json:"title" binding:"required,min=1,max=200"`
    Author string `json:"author" binding:"required,min=1,max=100"`
    Stock  int    `json:"stock" binding:"gte=0"`
}

// Handler 中
var dto BookCreateDTO
if err := c.ShouldBindJSON(&dto); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

### `ToResponse()` 方法

Model → DTO 轉換的標準做法：

```go
func (b *Book) ToResponse() BookResponseDTO {
    return BookResponseDTO{
        ID:     b.ID,
        Title:  b.Title,
        Author: b.Author,
        ISBN:   b.ISBN,
        Stock:  b.Stock,
    }
}
```

### 依賴注入

透過 constructor 傳入依賴，而非在 Service 內直接 new：

```go
// Good - 依賴注入
func NewBookService(repo repository.BookRepository) BookService {
    return &bookServiceImpl{repo: repo}
}

// Bad - 直接 new
func NewBookService() BookService {
    return &bookServiceImpl{repo: repository.NewBookRepository()}
}
```

---

## 動手做

### 必做

**1. 閱讀 Service 介面**

閱讀 `book_service.go` 介面，列出所有方法簽章並說明每個方法的用途。

**2. 追蹤 Create 流程**

追蹤 `book_service_impl.go` 的 `Create()` 完整流程：

```
DTO → Model → repo.Create → 回傳 ResponseDTO
```

**3. 重構 Handler**

把 Handler 中直接操作 Repository 的程式碼改為呼叫 Service：

```go
// 重構前
func (h *BookAPIHandler) GetBooks(c *gin.Context) {
    books, _ := h.repo.FindAll()
    c.JSON(http.StatusOK, books)
}

// 重構後
func (h *BookAPIHandler) GetBooks(c *gin.Context) {
    books, err := h.service.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, books)
}
```

確認頁面和 API 行為不變。

### 延伸挑戰

1. 在 Service 加入 `FindByISBN()` 方法（`⟵ 需完成 U06 必做 3`），串接 Repository 的新方法
2. 畫出 CreateDTO → Book → ResponseDTO 的資料流向圖

---

## 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `nil pointer dereference` | Service 或 Repository 沒有正確初始化 | 確認 constructor 有正確傳入依賴，並檢查初始化順序 |
| 介面方法簽章與實作不一致 | Go 是 implicit interface，簽章必須完全一致 | 編譯器會報錯，對照介面與實作的方法簽章修正 |

---

## 驗收標準（DoD）

- [ ] 能說明 Service 層存在的理由
- [ ] 能解釋 Model / CreateDTO / ResponseDTO 三者的職責差異
- [ ] 頁面和 API 行為與重構前完全一致
