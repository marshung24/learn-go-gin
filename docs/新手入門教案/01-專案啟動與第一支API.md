# U01｜專案啟動與第一支 API

> 能啟動專案並用瀏覽器呼叫自訂端點 ｜ 90 min ｜ 前置依賴：無

---

## 為什麼先教這個？

讓學員在第一堂課就看到「程式跑起來」的成就感。沒有什麼比親手讓一個 Web 應用回應你的請求更能建立學習信心。

---

## 對應檔案

| 檔案 | 角色 |
|------|------|
| `go.mod` / `go.sum` | 依賴管理（Go Modules） |
| `main.go` | 應用程式進入點 |
| `handler/hello.go` | 最簡易的 handler，示範如何回應 HTTP 請求 |

---

## 核心觀念

### 專案標準目錄架構

```
go-gin-sample/
├── go.mod / go.sum                       ← 依賴管理
├── docker-compose.yml                    ← 容器編排
├── main.go                               ← 進入點
├── internal/
│   ├── config/                           ← 設定載入
│   ├── handler/                          ← HTTP handler
│   ├── service/                          ← 業務邏輯
│   ├── repository/                       ← 資料存取
│   ├── model/                            ← Entity 與 DTO
│   ├── middleware/                       ← 中介軟體
│   └── cache/                            ← 快取相關
├── migrations/                           ← golang-migrate SQL
└── templates/                            ← html/template 模板
```

### Go Modules 做了什麼

`go.mod` 定義模組名稱與依賴版本，`go.sum` 記錄依賴的校驗碼：

```go
module github.com/example/go-gin-sample

go 1.22

require (
    github.com/gin-gonic/gin v1.10.0
    gorm.io/gorm v1.25.0
    // ...
)
```

常用指令：
- `go mod download` — 下載所有依賴
- `go mod tidy` — 清理未使用的依賴、補齊缺少的依賴

### `gin.Default()` vs `gin.New()`

```go
// Default 包含 Logger 和 Recovery middleware
r := gin.Default()

// New 是空白引擎，需要自己加 middleware
r := gin.New()
r.Use(gin.Logger(), gin.Recovery())
```

新手建議用 `Default()`，自動包含日誌與 panic 恢復。

### `c.JSON()` vs `c.String()`

```go
// 回傳 JSON（給 API 消費者）
c.JSON(http.StatusOK, gin.H{"message": "hello"})

// 回傳純文字
c.String(http.StatusOK, "Hello, Gin!")
```

### 內嵌 HTTP Server

不用另外裝 Web Server，`go run main.go` 就能跑，降低環境複雜度。

---

## 動手做

### 必做

**1. 啟動專案**

```bash
docker compose up -d         # 啟動 MySQL 與 Redis
go run main.go               # 啟動 Gin 應用
```

瀏覽器打開 `http://localhost:8080/hello`，確認看到回應。

**2. 新增自訂端點**

在 handler 新增一個 `/whoami` 端點，回傳自己的名字（JSON 格式）：

```go
func Whoami(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "name": "你的名字",
    })
}
```

用瀏覽器或 curl 測試：

```bash
curl http://localhost:8080/whoami
# 預期：{"name":"你的名字"}
```

### 延伸挑戰

1. 觀察 `go.mod` 的 require 區塊，寫下每個主要依賴的用途（如 `github.com/gin-gonic/gin` 提供什麼功能）

---

## 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `go run` 後報 port 衝突 | Port 8080 被其他程式占用 | macOS：`lsof -i :8080` 找出占用程式並關閉；或在程式碼中改 port |
| 修改 Go 檔後沒變化 | Go 不會自動重新編譯 | 需要重新執行 `go run main.go`，或使用 Air 進行 hot reload |
| `go: module not found` | 未執行 `go mod download` 或網路問題 | 執行 `go mod download`，確認 GOPROXY 設定正確 |

---

## 驗收標準（DoD）

- [ ] 能用瀏覽器或 `curl http://localhost:8080/whoami` 呼叫自訂端點並看到正確的 JSON 回應
- [ ] 能說出 `gin.Default()` 與 `gin.New()` 的差異
- [ ] 能說出 `go run main.go` 的作用
