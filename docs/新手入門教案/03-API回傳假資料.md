# U03｜API 回傳假資料

> 能用 Gin handler 回傳固定 JSON 並在 Swagger 測試 ｜ 60 min ｜ 前置依賴：U01

---

## 為什麼先教這個？

U02 做了給人看的頁面，這堂課做給機器用的 API。同樣用假資料，讓學員先理解「Gin Handler 怎麼回傳 JSON」。有了 API，就能用 Swagger UI 互動式測試，也為後面串接真實 DB 打下基礎。

---

## 對應檔案

| 檔案 | 角色 |
|------|------|
| `handler/book_api.go` | REST API handler（回傳 JSON），此刻先用假資料 |
| `docs/` | Swagger 文件目錄（swag init 產生） |

---

## 核心觀念

### `c.JSON()` 回傳 JSON

```go
func GetBooks(c *gin.Context) {
    books := []map[string]any{
        {"id": 1, "title": "Clean Code", "author": "Robert C. Martin"},
        {"id": 2, "title": "The Pragmatic Programmer", "author": "Andy Hunt"},
    }
    c.JSON(http.StatusOK, books)
}
```

`c.JSON()` 自動設定 `Content-Type: application/json`。

### 路由群組 `r.Group()`

```go
api := r.Group("/api")
{
    api.GET("/books", GetBooks)
    api.GET("/books/:id", GetBookByID)
    api.POST("/books", CreateBook)
    api.PUT("/books/:id", UpdateBook)
    api.DELETE("/books/:id", DeleteBook)
}
```

統一路徑前綴，REST API 慣例以 `/api/` 開頭。

### RESTful 動詞對應

| HTTP 動詞 | 用途 | 範例 |
|-----------|------|------|
| GET | 查詢資源 | `GET /api/books` |
| POST | 建立資源 | `POST /api/books` |
| PUT | 更新資源 | `PUT /api/books/1` |
| DELETE | 刪除資源 | `DELETE /api/books/1` |

### 路徑參數 vs Query String

```go
// 路徑參數 - /api/books/1
id := c.Param("id")

// Query String - /api/books?author=Martin
author := c.Query("author")
```

### swaggo/swag — Swagger 文件

用註解產生 Swagger 文件：

```go
// GetBooks godoc
// @Summary 取得所有書籍
// @Description 取得書籍清單
// @Tags books
// @Produce json
// @Success 200 {array} Book
// @Router /api/books [get]
func GetBooks(c *gin.Context) {
    // ...
}
```

執行 `swag init` 產生文件，訪問 `/swagger/index.html` 可互動式測試。

### 此刻不需要 DB

Handler 直接回傳假的 slice，先專注學 REST 概念。

---

## 動手做

### 必做

**1. 寫 GET /api/books**

寫一個 `GET /api/books` handler，回傳包含 3 本假書的 JSON 陣列：

```go
func GetBooks(c *gin.Context) {
    books := []gin.H{
        {"id": 1, "title": "Clean Code", "author": "Robert C. Martin", "stock": 5},
        {"id": 2, "title": "The Pragmatic Programmer", "author": "Andy Hunt", "stock": 3},
        {"id": 3, "title": "Refactoring", "author": "Martin Fowler", "stock": 0},
    }
    c.JSON(http.StatusOK, books)
}
```

**2. 寫 GET /api/books/:id**

寫一個 `GET /api/books/:id` handler，依 ID 回傳單筆假書（固定回傳同一本即可）：

```go
func GetBookByID(c *gin.Context) {
    id := c.Param("id")
    book := gin.H{
        "id":     id,
        "title":  "Clean Code",
        "author": "Robert C. Martin",
        "stock":  5,
    }
    c.JSON(http.StatusOK, book)
}
```

**3. 測試 Swagger UI**

執行 `swag init`，打開 Swagger UI（`http://localhost:8080/swagger/index.html`），測試剛才寫的 API。

### 延伸挑戰

1. 加入 `POST /api/books` handler，接收 JSON body 並用 `fmt.Println` 印出收到的資料（還沒有 DB，先確認能收到請求）

---

## 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| Swagger UI 看不到新加的 API | 沒有重新執行 `swag init` | 每次改 API 註解後都要重跑 `swag init` 重新產生文件 |
| `:id` 路由參數取不到 | 用錯方法取參數 | 使用 `c.Param("id")` 取路徑參數，`c.Query("key")` 取 query string |
| JSON 欄位名稱不對 | struct field 沒有 json tag | 加上 `json:"fieldName"` tag |

---

## 驗收標準（DoD）

- [ ] 能用 Swagger UI 或 `curl http://localhost:8080/api/books` 看到 JSON 回應
- [ ] 能說出路由群組 `r.Group()` 的用途
- [ ] 能說出 `c.Param()` 和 `c.Query()` 的差異
