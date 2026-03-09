# U08｜REST API 完善與錯誤處理

> 能用 middleware 統一錯誤回應，正確使用 HTTP 狀態碼 ｜ 90 min ｜ 前置依賴：U07

---

## 為什麼先教這個？

U07 完成了分層，但 API 還缺少正確的 HTTP 狀態碼和統一的錯誤回應。使用者送了錯誤資料、查了不存在的 ID，應該要拿到清楚的錯誤訊息而非一堆 stack trace。這堂課把 API 做到「生產品質」。

---

## 對應檔案

| 檔案 | 角色 |
|------|------|
| `handler/book_api.go` | REST API handler（加入正確的狀態碼與驗證） |
| `middleware/error_handler.go` | 全域錯誤處理 middleware |
| `middleware/logger.go` | Request logging middleware |

---

## 核心觀念

### HTTP 狀態碼語意

| 狀態碼 | 用途 | 範例 |
|--------|------|------|
| 200 OK | 請求成功 | GET 查詢成功 |
| 201 Created | 資源建立成功 | POST 新增成功 |
| 204 No Content | 成功但無回應內容 | DELETE 刪除成功 |
| 400 Bad Request | 請求格式錯誤 | JSON 欄位驗證失敗 |
| 404 Not Found | 資源不存在 | 查詢不存在的 ID |
| 500 Internal Server Error | 伺服器錯誤 | DB 連線失敗 |

### 正確的狀態碼回傳

```go
// 新增成功 — 201 Created
c.JSON(http.StatusCreated, book)

// 刪除成功 — 204 No Content
c.Status(http.StatusNoContent)

// 找不到 — 404 Not Found
c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
```

### Gin binding 驗證

```go
var dto BookCreateDTO
if err := c.ShouldBindJSON(&dto); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "error":   "Validation failed",
        "details": err.Error(),
    })
    return
}
```

`ShouldBindJSON` 自動觸發 struct tag 上的驗證規則，驗證失敗回 400。

### Recovery middleware

攔截所有 panic，避免整個 server 崩潰：

```go
r := gin.New()
r.Use(gin.Recovery())  // 攔截 panic，回傳 500
r.Use(gin.Logger())    // 記錄請求日誌
```

### 自訂錯誤類型

定義 `AppError` struct 統一錯誤格式：

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

// 常用錯誤
var (
    ErrNotFound      = &AppError{Code: 404, Message: "Resource not found"}
    ErrBadRequest    = &AppError{Code: 400, Message: "Invalid request"}
    ErrDuplicateISBN = &AppError{Code: 409, Message: "ISBN already exists"}
)
```

### 全域錯誤處理 Middleware

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            if appErr, ok := err.(*AppError); ok {
                c.JSON(appErr.Code, appErr)
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
            })
        }
    }
}
```

---

## 動手做

### 必做

**1. 修改 CreateBook — 回傳 201**

```go
func (h *BookAPIHandler) CreateBook(c *gin.Context) {
    var dto BookCreateDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    book, err := h.service.Create(&dto)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, book)  // 201 Created
}
```

**2. 修改 DeleteBook — 回傳 204**

```go
func (h *BookAPIHandler) DeleteBook(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

    if err := h.service.Delete(uint(id)); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    c.Status(http.StatusNoContent)  // 204 No Content
}
```

**3. 測試 400 錯誤**

故意送缺欄位的 JSON（如缺少 `title`），觀察 400 錯誤回應格式。

**4. 測試 404 錯誤**

查一筆不存在的 ID（如 `GET /api/books/99999`），觀察 404 錯誤回應。

### 延伸挑戰

1. 自訂一個業務錯誤（如 `ErrDuplicateISBN`），在 error handler middleware 中處理，回傳 409 Conflict

---

## 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| binding 驗證沒生效 | 用錯 `Bind` 方法 | 使用 `ShouldBindJSON` 而非 `BindJSON`（後者驗證失敗會自動回 400 但不可自訂格式） |
| 錯誤回應格式不一致 | 有些地方用 middleware，有些直接 c.JSON | 統一在 error handler middleware 處理所有錯誤 |

---

## 驗收標準（DoD）

- [ ] 能用 Swagger UI 完成完整 CRUD 操作（POST 201、DELETE 204）
- [ ] 能觸發 400 和 404 錯誤並觀察統一的錯誤回應格式
- [ ] 能說出 `ShouldBindJSON` 與 `BindJSON` 的差異
