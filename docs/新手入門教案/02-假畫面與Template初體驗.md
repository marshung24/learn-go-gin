# U02｜假畫面 — Template 初體驗

> 能用 html/template 做出書籍清單頁面（假資料） ｜ 60 min ｜ 前置依賴：U01

---

## 為什麼先教這個？

U01 讓程式跑起來了，但只有文字回應。這堂課讓學員第一次「看到畫面」——用 html/template 做出一個有書籍清單的 HTML 頁面。資料先用 Handler 中的假資料（hardcoded），不需要 DB，**重點是讓學員體驗「改了模板，重新整理就看到結果」的即時回饋感**。

> **講師提示**：先 demo 假資料版本，讓學員專心學 template 語法。後面 U06 串接 DB 時，只要改 Handler 的資料來源，頁面不用動——這就是模板引擎的價值。

---

## 對應檔案

| 檔案 | 角色 |
|------|------|
| `handler/book_view.go` | MVC handler（回傳 HTML 頁面），此刻先用假資料 |
| `templates/layout.html` | 共用版面（header / footer） |
| `templates/book/list.html` | 書籍清單頁 |
| `templates/book/detail.html` | 書籍詳情頁 |
| `templates/book/form.html` | 新增 / 編輯表單頁 |

---

## 核心觀念

### 伺服器端渲染（SSR）

html/template 在 server 端組好完整 HTML 再送出，不需要前端框架。

### `c.HTML()` 渲染模板

```go
func ListBooks(c *gin.Context) {
    books := []map[string]any{
        {"ID": 1, "Title": "Clean Code", "Author": "Robert C. Martin"},
        {"ID": 2, "Title": "The Pragmatic Programmer", "Author": "Andy Hunt"},
    }
    c.HTML(http.StatusOK, "book/list.html", gin.H{
        "books": books,
    })
}
```

- 第一個參數：HTTP 狀態碼
- 第二個參數：模板名稱
- 第三個參數：傳入模板的資料

### `gin.H{}`

等同於 `map[string]any{}`，方便快速建立傳給模板的資料。

### Template 關鍵語法

```html
<!-- 迴圈渲染 -->
{{range .books}}
<tr>
    <td>{{.ID}}</td>
    <td>{{.Title}}</td>
    <td>{{.Author}}</td>
</tr>
{{end}}

<!-- 輸出欄位值（自動 HTML escape 防 XSS） -->
<h1>{{.Title}}</h1>

<!-- 條件判斷 -->
{{if eq .Stock 0}}
<span class="text-danger">缺貨</span>
{{end}}

<!-- 套用共用版面 -->
{{template "layout" .}}
```

### 此刻不需要 DB

Handler 直接建立假的書籍清單丟給模板，先專注學 template 語法。

---

## 動手做

### 必做

**1. 觀察資料流**

閱讀 `book_view.go` 的 `ListBooks()` 如何把資料放進 `gin.H{}`，對照 `list.html` 中的 `{{range}}` 如何渲染每一行。

**2. 寫一個簡易版 Handler**

建立一個 Handler，回傳一個包含 3 本假書的清單，瀏覽 `http://localhost:8080/books` 看到表格：

```go
func ListBooks(c *gin.Context) {
    books := []map[string]any{
        {"ID": 1, "Title": "Clean Code", "Author": "Robert C. Martin", "Stock": 5},
        {"ID": 2, "Title": "The Pragmatic Programmer", "Author": "Andy Hunt", "Stock": 0},
        {"ID": 3, "Title": "Refactoring", "Author": "Martin Fowler", "Stock": 3},
    }
    c.HTML(http.StatusOK, "book/list.html", gin.H{
        "books": books,
    })
}
```

**3. 觀察詳情頁**

點進某本書的連結，觀察詳情頁如何用 `{{.Title}}` 顯示各欄位。

### 延伸挑戰

1. 修改 `list.html`，在表格加一欄「狀態」，當庫存 = 0 時顯示紅字「缺貨」（提示：`{{if eq .Stock 0}}`）

---

## 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `html/template: "xxx" is undefined` | 模板檔案未正確載入或名稱不符 | 確認 `LoadHTMLGlob()` 的 pattern 與模板檔案路徑一致 |
| 頁面顯示空白 | Handler 回傳了 JSON 而非 HTML | 確認使用 `c.HTML()` 而非 `c.JSON()` |
| 中文亂碼 | 模板檔案編碼非 UTF-8 | 確保模板檔案使用 UTF-8 編碼儲存 |

---

## 驗收標準（DoD）

- [ ] 能在瀏覽器看到一個有書籍清單的 HTML 頁面（即使資料是假的）
- [ ] 能說出 `{{range}}` 和 `{{.FieldName}}` 的用途
- [ ] 能說出 `c.HTML()` 三個參數的意義
