# 書籍借閱管理系統

> Go 1.22 / Gin 1.10 / GORM / golang-migrate / html/template / MySQL 8 / Redis 7 / Swagger (swaggo/swag)

以「書籍借閱管理」為案例的 Gin 練習專案，涵蓋常見分層架構：
書籍 CRUD（MySQL 持久化 + Redis 快取）、MVC 頁面與 REST API 並存、golang-migrate 資料庫版本管理、啟動預熱與 Redis 過期事件監聽。

## 前置條件

本專案搭配容器化開發環境使用（例如 GitHub Codespaces 或 Dev Container），
容器啟動後 MySQL 與 Redis 服務已就緒，不需另行安裝。

若要在本機直接執行，請自行準備以下服務：

- Go 1.22+
- MySQL 8.0（DB: `app`, User: `app`, Password: `app`）
  - DSN: `app:app@tcp(localhost:3306)/app?parseTime=true`
- Redis 7
  - Host: `localhost`, Port: `6379`

## 快速開始

### 使用 Dev Container（推薦）

1. 安裝 VS Code 擴充套件 **Dev Containers**（`ms-vscode-remote.remote-containers`）
2. Clone 專案並用 VS Code 開啟
3. 左下角點選 `><` 圖示 → 選擇 **Reopen in Container**
4. 等待容器建置完成（首次需較長時間），MySQL 與 Redis 會自動啟動並通過 healthcheck
5. 在容器內 Terminal 執行：

```bash
# 安裝依賴
go mod download

# 執行 migration
make migrate-up

# 啟動
go run main.go
```

> GitHub Codespaces 使用者：直接在 Codespaces 開啟專案即可，會自動套用 Dev Container 設定。

### 本機直接執行

請先確認[前置條件](#前置條件)中的服務已就緒，再執行：

```bash
git clone <repo-url> /workspaces/go-gin-sample
cd /workspaces/go-gin-sample

# 安裝依賴
go mod download

# 執行 migration
make migrate-up

# 啟動
go run main.go
```

啟動成功後驗證：

```bash
curl http://localhost:8080/hello          # Hello, Gin!
curl http://localhost:8080/api/books      # JSON 陣列，3 筆範例書籍
```

瀏覽器開啟 http://localhost:8080/books 可查看書籍清單頁。

## 可用網址

| 網址 | 說明 |
|------|------|
| http://localhost:8080/books | 書籍清單頁（html/template） |
| http://localhost:8080/books/{id} | 書籍詳情頁（html/template） |
| http://localhost:8080/swagger/index.html | Swagger UI — 互動式 API 文件與測試 |

> Swagger UI 僅顯示 `/api/books` 相關端點，MVC 頁面不會出現。

## API 端點

### REST API

| Method | Path | 說明 |
|--------|------|------|
| GET | `/api/books` | 取得所有書籍 |
| GET | `/api/books/{id}` | 取得單本書籍 |
| POST | `/api/books` | 新增書籍 |
| PUT | `/api/books/{id}` | 更新書籍 |
| DELETE | `/api/books/{id}` | 刪除書籍 |

### MVC 頁面

| Path | 說明 |
|------|------|
| `/books` | 書籍清單頁 |
| `/books/{id}` | 書籍詳情頁 |

## 專案結構

```
.devcontainer/
└── devcontainer.json              # Dev Container 設定（Codespaces / VS Code 遠端開發）
docker/
└── go/
    └── Dockerfile                 # Go 1.22 開發容器映像
docker-compose.yml                 # 三服務編排：app(bind mount) / db(MySQL 8) / redis(Redis 7)
.env.example                       # 環境變數範本（PORT、DB 連線等）
main.go                            # 應用程式進入點
internal/
├── config/
│   └── config.go                  # 設定載入（Viper）
├── handler/
│   ├── book_view.go               # MVC 頁面 handler
│   └── book_api.go                # REST API handler
├── service/
│   ├── book_service.go            # 介面定義
│   └── book_service_impl.go       # 業務邏輯（Cache-Aside）
├── repository/
│   ├── book_repository.go         # 介面定義
│   └── book_repository_impl.go    # GORM 實作
├── model/
│   ├── book.go                    # Entity
│   └── dto/
│       ├── book_create_dto.go     # 請求 DTO（binding 驗證）
│       └── book_response_dto.go   # 回應 DTO
├── middleware/
│   ├── error_handler.go           # 統一錯誤處理
│   └── logger.go                  # request logging
└── cache/
    ├── redis.go                   # Redis client 設定
    └── warmup.go                  # 啟動預熱

migrations/
├── 000001_create_books_table.up.sql
├── 000001_create_books_table.down.sql
├── 000002_insert_sample_data.up.sql
└── 000002_insert_sample_data.down.sql

templates/
├── layout.html
└── book/
    ├── list.html
    ├── detail.html
    └── form.html
```

## 設定檔與 Profile

| 設定檔 | DB Host | Redis Host | 適用環境 |
|--------|---------|------------|----------|
| `.env` | localhost | localhost | 本機直接開發 |
| `.env.docker` | db | redis | 容器環境（覆寫連線位址） |

容器環境已透過環境變數自動啟用 docker 設定，不需手動設定。
本機開發時使用預設的 localhost 連線。

### 容器服務與 Port 對應

| 服務 | 映像 | 容器內 Port | 預設對外 Port | 環境變數 |
|------|------|-------------|--------------|----------|
| app | `docker/go/Dockerfile` | 8080 | 8080 | `APP_PORT` |
| db | `mysql:8.0` | 3306 | 3306 | `MYSQL_PORT` |
| redis | `redis:7-alpine` | 6379 | 6379 | `REDIS_PORT` |

若需調整 Port 或專案名稱，複製 `.env.example` 為 `.env` 後修改即可：

```bash
cp .env.example .env
```

## 開發

### Hot Reload 模式

使用 Air 進行 hot reload：

```bash
# 安裝 Air
go install github.com/air-verse/air@latest

# 啟動 hot reload
air
```

### 測試

```bash
go test ./...                           # 執行全部測試
go test -v ./internal/service/...       # 只跑 service 測試
go test -cover ./...                    # 含覆蓋率
```

| 測試類 | 類型 | 說明 |
|--------|------|------|
| **book_service_test.go** | Mock 單元測試 | 使用 mockery 產生的 mock |
| **book_api_test.go** | HTTP handler 測試 | 使用 httptest |

## 教案文件

- [Go 新手入門教案（導覽）](docs/新手入門教案/README.md)
- [Go 新手入門教案綱要](docs/新手入門教案/Go新手入門教案綱要.md)
- [U01 — 專案啟動與第一支 API](docs/新手入門教案/01-專案啟動與第一支API.md)
- [U02 — 假畫面與 Template 初體驗](docs/新手入門教案/02-假畫面與Template初體驗.md)
- [U03 — API 回傳假資料](docs/新手入門教案/03-API回傳假資料.md)
- [U04 — 設定檔與多環境切換](docs/新手入門教案/04-設定檔與多環境切換.md)
- [U05 — golang-migrate 與資料庫設計](docs/新手入門教案/05-golang-migrate與資料庫設計.md)
- [U06 — GORM 串接資料庫](docs/新手入門教案/06-GORM串接資料庫.md)
- [U07 — 正規化架構與 Service 分層](docs/新手入門教案/07-正規化架構與Service分層.md)
- [U08 — REST API 完善與錯誤處理](docs/新手入門教案/08-REST-API完善與錯誤處理.md)
- [U09 — Redis 快取整合](docs/新手入門教案/09-Redis快取整合.md)
- [U10 — 測試策略與實戰](docs/新手入門教案/10-測試策略與實戰.md)
