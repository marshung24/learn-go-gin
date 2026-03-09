# Go 新手入門教案綱要

> 以本 Repo「書籍管理系統」為實作主軸，從零帶出 Go + Gin 全棧開發核心技能。
>
> **給教案作者的提醒**：本綱要為「骨架級」指導文件，內容刻意詳盡以確保教案撰寫時不遺漏關鍵要素。依此綱要寫出的**教案本身應力求精練**——每個單元的講義以「此刻實作所需的最小可行知識」為上限，避免資訊過載。
>
> **單元排列原則**：遵循「對應專案流程的漸進式推演」——教案對應專案推進流程（需求分析 → 討論 → 拆解 → 規劃 → 設計 → 開發 → 測試），在開發階段以「可視性」為脈絡，先讓學員看到畫面（假資料），再串接真實資料庫完成 CRUD 循環，最後才抽出 Service/DTO 分層與進階功能。每一步都有可見的產出，避免「寫了半天什麼都看不到」的挫折。

---

## §1 教案資訊

| 項目 | 說明 |
|------|------|
| 對象 | 有基礎程式概念（變數、迴圈、函式）與基本命令列操作經驗，但**尚未做過 Go Web 專案**的初學者 |
| 目標 | 獨立完成一個「能跑、能測、能部署」的 CRUD Web 應用（書籍管理系統）；所有單元的成果最終匯聚為同一個可動可用的範例網站 |
| 技術棧 | Go 1.22、Gin 1.10、GORM、golang-migrate、MySQL 8、Redis 7、html/template、swaggo/swag（Swagger）、Viper（設定管理）、Docker Compose |
| 時數 | 課前準備（自學）+ 10 單元（各 60–120 分鐘，詳見 §5 地圖層） |
| 教學方式 | 觀念講解 → Repo 實作示範 → 學員動手做 → Code Review → 驗收 |
| 先備知識 | 基礎 SQL（SELECT / INSERT）、HTTP 概念（GET / POST）、Git 基本操作（clone / commit） |
| 不含範圍 | 前端框架（React / Vue）、JWT 認證、gRPC、微服務架構、Kubernetes、Context 進階用法 |

---

## §2 學習成果（結訓時能做到）

1. **能說明** Gin 專案分層架構（Router/Handler → Service → Repository → Model）各層的職責與呼叫方向，理解 Go 的 package 與 interface 管理分層的設計哲學
2. **能使用** Go Modules 建置專案、用 Docker Compose 一鍵啟動整套開發環境（App + MySQL + Redis）
3. **能撰寫** golang-migrate migration 腳本管理資料庫 schema 演進，理解 up/down migration 與版本編號策略
4. **能實作** GORM Repository 完成 CRUD 操作，新增自訂查詢方法
5. **能設計** RESTful API（正確使用 HTTP 動詞與狀態碼），用 middleware 統一錯誤回應與 request logging，並產出 Swagger 互動式文件
6. **能製作** html/template 伺服器端渲染頁面（列表、表單、詳情），理解 PRG 模式
7. **能實作** Redis Cache-Aside 快取模式，並用 redis-cli 驗證快取行為
8. **能撰寫** Go testing 單元測試與 handler 測試，使用 mock interface 驗證依賴替身與基本覆蓋率

---

## §3 閱讀指引

### 練習標記說明

| 標記 | 意義 | 使用場景 |
|------|------|----------|
| **必做** | 課堂內所有學員都應完成的核心練習 | 達成該單元最低驗收門檻 |
| **延伸挑戰** | 給進度快的學員或課後作業 | 加深理解、跨單元整合 |

### 依賴標記說明

- 標有 `⟵ 需完成 UXX` 的項目，表示需先完成指定單元的對應練習才能進行
- 標有 `⟵ 需完成 UXX 延伸` 的項目，特指依賴某單元的「延伸挑戰」成果

### 命名規範

| 項目 | 格式 | 範例 |
|------|------|------|
| 單元編號 | U00–U10 | U03 |
| 檔案路徑 | 省略 `internal/`，保持可讀性 | `handler/book_api.go` |
| Migration 路徑 | 完整路徑 | `migrations/000001_create_books_table.up.sql` |
| 測試路徑 | 省略 `internal/` | `service/book_service_test.go` |

### 單元內部結構

每個單元依序包含六個子區塊：① 為什麼先教這個？ → ② 對應檔案 → ③ 核心觀念 → ④ 動手做 → ⑤ 踩坑提示 → ⑥ 驗收標準（DoD）

---

## §4 課前準備（U00）

> 本單元為課前自修，不佔正課時間。目標：確保第一堂課不卡在環境問題上。

### 安裝清單

| 工具 | 版本 | 推薦安裝方式（macOS） | Windows WSL 注意事項 |
|------|------|----------------------|---------------------|
| Go | 1.22+ | [官方安裝](https://go.dev/dl/) 或 `brew install go` | WSL2 中同樣使用官方安裝或 apt |
| Docker Desktop | 最新版 | [Docker Desktop for Mac](https://www.docker.com/products/docker-desktop/) | 安裝 [Docker Desktop for Windows](https://docs.docker.com/desktop/setup/install/windows-install/) 並啟用 WSL2 後端整合 |
| Git | 最新版 | `brew install git` 或 Xcode Command Line Tools | WSL2 內建或 `sudo apt install git` |
| IDE（主要） | — | **VS Code** + [Go 擴充套件](https://marketplace.visualstudio.com/items?itemName=golang.Go) | 同 macOS；可搭配 [Remote - WSL 擴充](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-wsl) |
| IDE（替代） | — | GoLand（參考 [官方安裝指南](https://www.jetbrains.com/go/download/)） | 同 macOS |

> **IDE 選擇說明**：本教案以 **VS Code** 為主要示範環境。GoLand 功能更豐富但需付費，有興趣的學員可參考 [GoLand Getting Started](https://www.jetbrains.com/help/go/getting-started.html) 自行探索。兩者在本教案範圍內的操作差異極小。

> **作業系統說明**：本教案以 **macOS** 為主要示範環境。Windows 使用者建議透過 **WSL2**（Windows Subsystem for Linux）進行開發，可獲得與 macOS/Linux 一致的命令列體驗。

### 驗證指令

```bash
go version                       # 確認顯示 go1.22.x 或更高
docker compose version           # 確認 Docker Compose 可用
git clone <本 Repo URL> && cd go-gin-sample
go mod download                  # 確認能下載依賴
```

### 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `docker: command not found` | Docker Desktop 未啟動或未安裝 | macOS / Windows 需先啟動 Docker Desktop 應用程式 |
| `go mod download` 卡住或逾時 | 公司網路可能封鎖 proxy.golang.org | 設定 `GOPROXY=https://goproxy.io,direct` 或切換至個人網路 |
| `go version` 顯示舊版本 | 系統有多版本 Go | 確認 PATH 中 Go 1.22 路徑在前，或用 `go1.22.x download` 安裝指定版本 |

---

## §5 課程模組

### 地圖層（單元總覽）

| 單元 | 名稱 | 時數 | 一句話目標 | 前置依賴 |
|------|------|------|-----------|----------|
| U01 | 專案啟動與第一支 API | 90 min | 能啟動專案並用瀏覽器呼叫自訂端點 | — |
| U02 | 假畫面 — Template 初體驗 | 60 min | 能用 html/template 做出一個書籍清單頁面（假資料） | U01 |
| U03 | API 回傳假資料 | 60 min | 能用 Gin handler 回傳固定 JSON 並在 Swagger 測試 | U01 |
| U04 | 設定檔與多環境切換 | 60 min | 能解釋環境變數機制並用 Docker Compose 啟動完整環境 | U01 |
| U05 | golang-migrate 與資料庫設計 | 90 min | 能撰寫 migration SQL 並執行 up/down | U04 |
| U06 | GORM 串接資料庫 | 120 min | 能用 GORM 完成真實 CRUD，頁面與 API 顯示 DB 真實資料 | U02, U03, U05 |
| U07 | 正規化架構 — Service 與 DTO 分層 | 90 min | 能從 Handler 抽出 Service 層，引入 DTO 做資料驗證與格式轉換 | U06 |
| U08 | REST API 完善與錯誤處理 | 90 min | 能用 middleware 統一錯誤回應，正確使用 HTTP 狀態碼 | U07 |
| U09 | Redis 快取整合 | 90 min | 能用 redis-cli 驗證快取行為並說明 Cache-Aside 流程 | U07 |
| U10 | 測試策略與實戰 | 120 min | 能為新方法撰寫單元測試與 handler 測試 | U08 |

> **設計說明**：單元順序遵循「對應專案流程的漸進式推演」——對應開發階段的可視性脈絡：環境確認（U01）→ 假畫面 Level-1 UI（U02）→ 假資料 Level-2 程式（U03）→ 設定環境（U04）→ 資料庫設計 Level-3 資料（U05）→ 串接真實 DB（U06）→ 架構正規化（U07）→ 完善 API（U08）→ 效能優化（U09）→ 測試驗證（U10）。學員在 U02 就能在瀏覽器看到書籍清單頁面，U06 結束時已有一個從 DB 讀寫的完整 CRUD 網站，之後才逐步重構為正規架構。

---

### U01｜專案啟動與第一支 API（90 min）

#### ① 為什麼先教這個？

讓學員在第一堂課就看到「程式跑起來」的成就感。沒有什麼比親手讓一個 Web 應用回應你的請求更能建立學習信心。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `go.mod` / `go.sum` | 依賴管理（Go Modules） |
| `main.go` | 應用程式進入點 |
| `handler/hello.go` | 最簡易的 handler，示範如何回應 HTTP 請求 |

#### ③ 核心觀念

- **專案標準目錄架構** — `cmd/`（進入點）、`internal/`（私有程式碼）、`pkg/`（公開套件），Go 社群慣例的分層方式
- **Go Modules 做了什麼** — 管理依賴版本、自動下載套件；`go mod download` 和 `go mod tidy` 是常用指令
- **`gin.Default()` vs `gin.New()`** — Default 包含 Logger 和 Recovery middleware，New 是空白引擎
- **`c.JSON()` vs `c.String()`** — 前者回傳 JSON 給 API 消費者，後者回傳純文字
- **內嵌 HTTP Server** — 不用另外裝 Web Server，`go run main.go` 就能跑，降低環境複雜度

#### ④ 動手做

**必做**
1. 執行 `docker compose up -d` 啟動 MySQL 與 Redis，再執行 `go run main.go`，瀏覽器打開 `http://localhost:8080/hello`，確認看到回應
2. 在 handler 新增一個 `/whoami` 端點，回傳自己的名字（JSON 格式）

**延伸挑戰**
1. 觀察 `go.mod` 的 require 區塊，寫下每個主要依賴的用途（如 `github.com/gin-gonic/gin` 提供什麼功能）

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `go run` 後報 port 衝突 | Port 8080 被其他程式占用 | macOS：`lsof -i :8080` 找出占用程式並關閉；或在程式碼中改 port |
| 修改 Go 檔後沒變化 | Go 不會自動重新編譯 | 需要重新執行 `go run main.go`，或使用 Air 進行 hot reload |
| `go: module not found` | 未執行 `go mod download` 或網路問題 | 執行 `go mod download`，確認 GOPROXY 設定正確 |

#### ⑥ 驗收標準（DoD）

能用瀏覽器或 `curl http://localhost:8080/whoami` 呼叫自訂端點並看到正確的 JSON 回應。

---

### U02｜假畫面 — Template 初體驗（60 min）

#### ① 為什麼先教這個？

U01 讓程式跑起來了，但只有文字回應。這堂課讓學員第一次「看到畫面」——用 html/template 做出一個有書籍清單的 HTML 頁面。資料先用 Handler 中的假資料（hardcoded），不需要 DB，**重點是讓學員體驗「改了模板，重新整理就看到結果」的即時回饋感**。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `handler/book_view.go` | MVC handler（回傳 HTML 頁面），此刻先用假資料 |
| `templates/layout.html` | 共用版面（header / footer） |
| `templates/book/list.html` | 書籍清單頁 |
| `templates/book/detail.html` | 書籍詳情頁 |
| `templates/book/form.html` | 新增 / 編輯表單頁 |

#### ③ 核心觀念

- **伺服器端渲染（SSR）** — html/template 在 server 端組好完整 HTML 再送出，不需要前端框架
- **`c.HTML()` 渲染模板** — 第一個參數是 HTTP 狀態碼，第二個是模板名稱，第三個是傳入模板的資料
- **`gin.H{}`** — 等同於 `map[string]any{}`，方便快速建立傳給模板的資料
- **Template 關鍵語法**：
  - `{{range .Items}}...{{end}}` — 迴圈渲染
  - `{{.FieldName}}` — 輸出欄位值（自動 HTML escape 防 XSS）
  - `{{template "layout" .}}` — 套用共用版面
- **此刻不需要 DB** — Handler 直接建立假的書籍清單丟給模板，先專注學 template 語法

> **講師提示**：先 demo 假資料版本，讓學員專心學 template 語法。後面 U06 串接 DB 時，只要改 Handler 的資料來源，頁面不用動——這就是模板引擎的價值。

#### ④ 動手做

**必做**
1. 觀察 `book_view.go` 的 `ListBooks()` 如何把資料放進 `gin.H{}`，對照 `list.html` 中的 `{{range}}` 如何渲染每一行
2. 寫一個簡易版 Handler，回傳一個包含 3 本假書的清單，瀏覽 `http://localhost:8080/books` 看到表格
3. 點進某本書的連結，觀察詳情頁如何用 `{{.Title}}` 顯示各欄位

**延伸挑戰**
1. 修改 `list.html`，在表格加一欄「狀態」，當庫存 = 0 時顯示紅字「缺貨」（提示：`{{if eq .Stock 0}}`）

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `html/template: "xxx" is undefined` | 模板檔案未正確載入或名稱不符 | 確認 `LoadHTMLGlob()` 的 pattern 與模板檔案路徑一致 |
| 頁面顯示空白 | Handler 回傳了 JSON 而非 HTML | 確認使用 `c.HTML()` 而非 `c.JSON()` |
| 中文亂碼 | 模板檔案編碼非 UTF-8 | 確保模板檔案使用 UTF-8 編碼儲存 |

#### ⑥ 驗收標準（DoD）

能在瀏覽器看到一個有書籍清單的 HTML 頁面（即使資料是假的），並能說出 `{{range}}` 和 `{{.FieldName}}` 的用途。

---

### U03｜API 回傳假資料（60 min）

#### ① 為什麼先教這個？

U02 做了給人看的頁面，這堂課做給機器用的 API。同樣用假資料，讓學員先理解「Gin Handler 怎麼回傳 JSON」。有了 API，就能用 Swagger UI 互動式測試，也為後面串接真實 DB 打下基礎。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `handler/book_api.go` | REST API handler（回傳 JSON），此刻先用假資料 |
| `docs/` | Swagger 文件目錄（swag init 產生） |

#### ③ 核心觀念

- **`c.JSON(http.StatusOK, data)`** — 回傳 JSON，自動設定 Content-Type
- **路由群組 `r.Group("/api")`** — 統一路徑前綴，REST API 慣例以 `/api/` 開頭
- **RESTful 動詞對應** — GET（查）、POST（建）、PUT（改）、DELETE（刪）
- **swaggo/swag** — 用註解產生 Swagger 文件，訪問 `/swagger/index.html` 可互動式測試
- **此刻不需要 DB** — Handler 直接回傳假的 slice，先專注學 REST 概念

#### ④ 動手做

**必做**
1. 寫一個 `GET /api/books` handler，回傳包含 3 本假書的 JSON 陣列
2. 寫一個 `GET /api/books/:id` handler，依 ID 回傳單筆假書（固定回傳同一本即可）
3. 執行 `swag init`，打開 Swagger UI（`http://localhost:8080/swagger/index.html`），測試剛才寫的 API

**延伸挑戰**
1. 加入 `POST /api/books` handler，接收 JSON body 並用 `fmt.Println` 印出收到的資料（還沒有 DB，先確認能收到請求）

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| Swagger UI 看不到新加的 API | 沒有重新執行 `swag init` | 每次改 API 註解後都要重跑 `swag init` 重新產生文件 |
| `:id` 路由參數取不到 | 用錯方法取參數 | 使用 `c.Param("id")` 取路徑參數，`c.Query("key")` 取 query string |
| JSON 欄位名稱不對 | struct field 沒有 json tag | 加上 `json:"fieldName"` tag |

#### ⑥ 驗收標準（DoD）

能用 Swagger UI 或 `curl http://localhost:8080/api/books` 看到 JSON 回應，並能說出路由群組和路徑參數的用途。

---

### U04｜設定檔與多環境切換（60 min）

#### ① 為什麼先教這個？

接下來要串接真實 DB（U05-U06），得先搞懂「DB 連線設定在哪、怎麼改、怎麼切環境」。不懂設定管理，後面 MySQL 連不上就會卡住。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `config/config.go` | Viper 設定載入邏輯 |
| `.env.example` | 環境變數範本 |
| `docker-compose.yml` | 多容器編排（App + MySQL + Redis） |

#### ③ 核心觀念

- **Viper 設定管理** — 支援 .env 檔、環境變數、YAML 等多種來源，自動綁定 struct
- **環境變數優先於檔案** — 容器環境可用環境變數覆寫設定，不需改檔案
- **`.env` 檔案** — 開發環境常用的設定方式，不應 commit 到 git（用 `.env.example` 作範本）
- **Docker Compose 一鍵啟動** — 一個 `docker compose up -d` 拉起 App + MySQL + Redis 三個容器

> **開發方式補充**：本教案以 **Docker Compose** 為主要開發環境方案。其他常見方式包括本機直接安裝 MySQL + Redis、Podman、DevContainer 等。

#### ④ 動手做

**必做**
1. 執行 `docker compose up -d`，用 `docker compose ps` 確認三個容器（app / db / redis）都在 running 狀態
2. 比對 `.env.example` 中本機與容器環境的設定差異，說明為什麼容器內 DB host 是 `db` 而非 `localhost`

**延伸挑戰**
1. 新增一個自訂設定項（如 `APP_NAME`），在 config struct 中讀取並在啟動時印出

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `docker compose up` 後 App 啟動失敗 | MySQL 容器尚未 ready | 等 healthcheck 通過後重試，或查看 `docker compose logs app` 確認錯誤 |
| `.env` 變數沒生效 | 忘了從 `.env.example` 複製為 `.env` | `cp .env.example .env` 後重新啟動 |
| Viper 讀不到環境變數 | 沒有呼叫 `viper.AutomaticEnv()` | 確認 config 初始化時有啟用環境變數自動綁定 |

#### ⑥ 驗收標準（DoD）

能解釋本機與容器兩套設定的差異，並能成功用 Docker Compose 啟動完整環境。

---

### U05｜golang-migrate 與資料庫設計（90 min）

#### ① 為什麼先教這個？

U02-U03 用的是假資料，接下來要換成真的。第一步是「設計資料庫結構」——用 golang-migrate 管理 schema 版本，確保多人協作時每個人的 DB 長一樣。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `migrations/000001_create_books_table.up.sql` | 建表（初始 schema） |
| `migrations/000001_create_books_table.down.sql` | 回滾建表 |
| `migrations/000002_insert_sample_data.up.sql` | 塞入範例資料（seed data） |
| `migrations/000002_insert_sample_data.down.sql` | 回滾範例資料 |
| `Makefile` | migration 指令封裝 |

#### ③ 核心觀念

- **命名慣例** — `{序號}_{描述}.{up|down}.sql`，序號遞增；up 是正向 migration，down 是回滾
- **`schema_migrations` 表** — golang-migrate 在 DB 中自動建立的紀錄表，記錄當前版本
- **up/down 對稱** — 每個 up.sql 都要有對應的 down.sql，確保能回滾
- **欄位設計思考** — 為什麼用 `BIGINT` 不用 `INT`？為什麼 `isbn` 加 `UNIQUE`？`DEFAULT CURRENT_TIMESTAMP` 與 `ON UPDATE CURRENT_TIMESTAMP` 的差異
- **`migrate` CLI** — `migrate -path migrations -database "mysql://..." up/down/version`

#### ④ 動手做

**必做**
1. 閱讀 000001 和 000002 的 SQL，說出每個欄位的設計理由
2. 新增 `000003_add_publisher_column.up.sql` 和 `.down.sql`，為 books 表加一個 `publisher VARCHAR(100)` 欄位
3. 執行 `make migrate-up`，確認 migration 成功；連線 MySQL 查詢 `SELECT * FROM schema_migrations` 確認版本

**延伸挑戰**
1. 執行 `make migrate-down` 回滾一版，確認 publisher 欄位消失；思考：什麼情況下會需要 down migration？

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `dirty database version` | 上次 migration 執行失敗但版本號已更新 | 手動修正 `schema_migrations` 表的 dirty 欄位為 false，或 `migrate force VERSION` |
| migration 檔案被忽略 | 序號不連續或格式錯誤 | 確認序號遞增且格式為 `{序號}_{描述}.{up|down}.sql` |

#### ⑥ 驗收標準（DoD）

能撰寫一組新的 migration SQL（up + down）並成功執行，且能在 `schema_migrations` 確認版本紀錄。

---

### U06｜GORM 串接資料庫（120 min）

#### ① 為什麼先教這個？

U05 設計好了資料表，這堂課要把假資料換成真的——用 GORM 串接 MySQL，讓 Handler 直接從 DB 讀寫資料。**完成後學員會第一次看到「頁面上的資料來自真實資料庫」的完整 CRUD 循環。** 此刻先不抽 Service 層，Handler 直接操作 Repository——保持最短路徑讓 CRUD 跑起來。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `model/book.go` | DB Model（對應 books 表） |
| `repository/book_repository.go` | Repository 介面定義 |
| `repository/book_repository_impl.go` | GORM 實作 |

#### ③ 核心觀念

- **GORM Model 對應** — struct tag `gorm:"column:xxx"` 對應 DB 欄位；`gorm:"primaryKey"` 標示主鍵
- **GORM 自動遷移** — `db.AutoMigrate(&Book{})` 可自動同步 schema（教學用，生產環境建議用 migration）
- **Repository Pattern** — 定義介面，實作與介面分離，方便測試時 mock
- **`db.First()` vs `db.Find()`** — First 取一筆（找不到會 error），Find 取多筆（找不到回傳空 slice）
- **替換假資料** — 修改 U02/U03 的 Handler，把 hardcoded slice 改成呼叫 repository，頁面立刻顯示真實 DB 資料

> **講師提示**：這是課程的第一個「啊哈時刻」——學員親眼看到頁面從假資料切換到真實 DB 資料。建議 live demo 時先秀假資料頁面，再切換到 Repository 呼叫，重新整理瀏覽器的瞬間最有衝擊力。

#### ④ 動手做

**必做**
1. 閱讀 `book_repository.go` 介面和 `book_repository_impl.go` 實作，理解每個方法對應的 GORM 操作
2. 修改 `book_view.go`，把 `ListBooks()` 方法中的假資料替換為 `repo.FindAll()`，重新整理頁面確認顯示 DB 資料
3. 在 Repository 新增 `FindByISBN(isbn string) (*Book, error)` 方法，確認編譯通過

**延伸挑戰**
1. 新增一個多條件查詢方法（如依作者 + 庫存範圍），體驗 GORM 的 `Where` 鏈式呼叫

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `record not found` panic | 用 `First()` 查不存在的資料時會 error | 檢查 `errors.Is(err, gorm.ErrRecordNotFound)` 或改用 `Find()` |
| 欄位值全是零值 | Model struct 的欄位名稱與 DB column 不對應 | 加上正確的 `gorm:"column:xxx"` tag |
| 連線失敗 | DSN 格式錯誤或 DB 服務未啟動 | 確認 docker compose 服務狀態，檢查 DSN 格式 |

#### ⑥ 驗收標準（DoD）

能在瀏覽器看到來自 DB 的書籍清單（不再是假資料），並能新增一個 Repository 方法。

---

### U07｜正規化架構 — Service 與 DTO 分層（90 min）

#### ① 為什麼先教這個？

U06 完成了 CRUD，但 Handler 直接操作 Repository——所有邏輯擠在一起，改一個地方可能影響全部。這堂課把程式碼「整理乾淨」：抽出 Service 層集中業務邏輯，引入 DTO 分離輸入/輸出格式。**重構前後功能完全一樣，但程式碼更好維護、更好測試。**

> **講師提示**：先展示重構前（Handler 直接操作 Repository）和重構後（Handler → Service → Repository）的對比，讓學員感受到「程式碼可以一樣能跑，但結構差很多」。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `service/book_service.go` | Service 介面（定義業務合約） |
| `service/book_service_impl.go` | Service 實作（業務邏輯，此刻先不含快取） |
| `model/dto/book_create_dto.go` | 輸入驗證 DTO（新增/更新用） |
| `model/dto/book_response_dto.go` | 輸出格式 DTO（API 回傳用） |

#### ③ 核心觀念

- **介面與實作分離** — Go 的 implicit interface 讓測試可以用 Mock 替換實作，也方便日後抽換不同實作
- **業務責任邊界** — Service 負責「業務規則」（如檢查重複 ISBN），Handler 負責「HTTP 協議」（如回傳 201 狀態碼），兩者不應互相滲透
- **Model 與 DTO 分離** — `Book` 對應 DB，`BookCreateDTO` 對應「收什麼」，`BookResponseDTO` 對應「給什麼」。修改其中一邊不會牽動另一邊
- **Gin binding 驗證** — `binding:"required,min=1"` 在 DTO struct tag 上宣告驗證規則，搭配 `c.ShouldBindJSON()` 自動觸發
- **`ToResponse()` 方法** — Model → DTO 轉換的標準做法，轉換邏輯集中在一處
- **依賴注入** — 透過 constructor 傳入依賴（如 Repository），而非在 Service 內直接 new

#### ④ 動手做

**必做**
1. 閱讀 `book_service.go` 介面，列出所有方法簽章並說明每個方法的用途
2. 追蹤 `book_service_impl.go` 的 `Create()` 完整流程：DTO → Model → repo.Create → 回傳 ResponseDTO
3. 重構 Handler：把直接操作 Repository 的程式碼改為呼叫 Service，確認頁面和 API 行為不變

**延伸挑戰**
1. 在 Service 加入 `FindByISBN()` 方法（`⟵ 需完成 U06 必做 3`），串接 Repository 的新方法
2. 畫出 CreateDTO → Book → ResponseDTO 的資料流向圖

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `nil pointer dereference` | Service 或 Repository 沒有正確初始化 | 確認 constructor 有正確傳入依賴，並檢查初始化順序 |
| 介面方法簽章與實作不一致 | Go 是 implicit interface，簽章必須完全一致 | 編譯器會報錯，對照介面與實作的方法簽章修正 |

#### ⑥ 驗收標準（DoD）

能說明 Service 層存在的理由，能解釋 Model/CreateDTO/ResponseDTO 三者的職責差異，頁面和 API 行為與重構前完全一致。

---

### U08｜REST API 完善與錯誤處理（90 min）

#### ① 為什麼先教這個？

U07 完成了分層，但 API 還缺少正確的 HTTP 狀態碼和統一的錯誤回應。使用者送了錯誤資料、查了不存在的 ID，應該要拿到清楚的錯誤訊息而非一堆 stack trace。這堂課把 API 做到「生產品質」。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `handler/book_api.go` | REST API handler（加入正確的狀態碼與驗證） |
| `middleware/error_handler.go` | 全域錯誤處理 middleware |
| `middleware/logger.go` | Request logging middleware |

#### ③ 核心觀念

- **HTTP 狀態碼語意** — 200 OK、201 Created、204 No Content、400 Bad Request、404 Not Found、500 Internal Server Error
- **`c.JSON(http.StatusCreated, data)`** — 新增成功回 201 而非預設的 200
- **Gin binding 驗證** — `c.ShouldBindJSON(&dto)` 自動觸發 struct tag 上的驗證，驗證失敗回 400
- **Recovery middleware** — 攔截所有 panic，避免整個 server 崩潰
- **自訂錯誤類型** — 定義 `AppError` struct 統一錯誤格式（code + message）

#### ④ 動手做

**必做**
1. 修改 `CreateBook` handler，成功時回傳 201、驗證失敗時回傳 400
2. 修改 `DeleteBook` handler，成功時回傳 204（No Content）
3. 故意送缺欄位的 JSON（如缺少 `title`），觀察 400 錯誤回應格式
4. 查一筆不存在的 ID（如 `GET /api/books/99999`），觀察 404 錯誤回應

**延伸挑戰**
1. 自訂一個業務錯誤（如 `ErrDuplicateISBN`），在 error handler middleware 中處理，回傳 409 Conflict

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| binding 驗證沒生效 | 用錯 `Bind` 方法 | 使用 `ShouldBindJSON` 而非 `BindJSON`（後者驗證失敗會自動回 400 但不可自訂格式） |
| 錯誤回應格式不一致 | 有些地方用 middleware，有些直接 c.JSON | 統一在 error handler middleware 處理所有錯誤 |

#### ⑥ 驗收標準（DoD）

能用 Swagger UI 完成完整 CRUD 操作（POST 201、DELETE 204），並能觸發 400 和 404 錯誤觀察統一的錯誤回應格式。

---

### U09｜Redis 快取整合（90 min）

#### ① 為什麼先教這個？

資料庫查詢是 Web 應用最慢的環節。把「熱資料」放在 Redis 記憶體中，回應速度可以從毫秒降到微秒等級。這堂課在 U07 建立的 Service 層上加入快取邏輯。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `service/book_service_impl.go` | Cache-Aside 讀寫邏輯的實際所在 |
| `cache/redis.go` | Redis client 初始化與設定 |
| `cache/warmup.go` | 啟動時快取預熱（preload 常用資料） |

#### ③ 核心觀念

- **Cache-Aside 模式** — 讀：查快取 → miss → 查 DB → 寫快取；寫/刪：操作 DB → 刪除快取
- **go-redis 基本操作** — `rdb.Get(ctx, key)` / `rdb.Set(ctx, key, value, ttl)` / `rdb.Del(ctx, key)`
- **JSON 序列化** — 用 `json.Marshal()` / `json.Unmarshal()` 將 struct 存入 Redis
- **快取 key 設計** — `book:{id}`，TTL 10 分鐘
- **為什麼 update / delete 要清快取？** — 避免讀到過期資料（快取一致性）
- **預熱（Warmup）** — 應用啟動時預先載入常用資料到快取，減少 cold start 的 cache miss

#### ④ 動手做

**必做**
1. 追蹤 `book_service_impl.go` 的 `FindByID()` 完整流程：Redis hit → 直接回傳 / Redis miss → 查 DB → 寫回 Redis
2. 用 `redis-cli` 觀察啟動後的快取 key：`KEYS book:*`（教學環境限定指令）
3. 手動刪除一個 key（`DEL book:1`），再呼叫 API 觀察 console log 中的 cache miss 與回寫

**延伸挑戰**
1. 呼叫 `GET /api/books/1` 兩次，比較回應時間差異；將 TTL 改為 1 分鐘觀察快取命中率變化
2. 思考：什麼情況下應該用 Cache-Aside？什麼情況下不適合？

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `KEYS` 指令在正式環境不能用 | `KEYS` 會阻塞 Redis（全表掃描） | 教學環境直觀好用；正式環境改用 `SCAN` |
| 快取與 DB 資料不一致 | 更新 DB 後忘了清快取 | 確保 Update / Delete 操作都有對應的 `rdb.Del(ctx, key)` |
| Redis 連線失敗 | 容器未啟動或設定錯誤 | 檢查 docker compose 狀態，確認 Redis host/port 設定 |

#### ⑥ 驗收標準（DoD）

能用 `redis-cli` 驗證快取行為（確認 key 存在、手動刪除後觀察回源），並能說明 Cache-Aside 的讀寫流程與快取失效策略。

---

### U10｜測試策略與實戰（120 min）

#### ① 為什麼先教這個？

不寫測試的程式碼是「可能會動」，有測試的程式碼是「確定會動」。測試是修改程式碼後不會造成回歸錯誤的安全網——這是從「學習專案」走向「生產專案」的關鍵一步。

#### ② 對應檔案

| 檔案 | 角色 |
|------|------|
| `service/book_service_test.go` | Service 單元測試（用 mock 隔離依賴） |
| `handler/book_api_test.go` | Handler 測試（用 httptest） |
| `mocks/` | mockery 產生的 mock 檔案 |

#### ③ 核心觀念

- **測試金字塔** — 單元測試（多且快）→ 整合測試（少且慢），先追求單元測試覆蓋
- **Go testing 基礎**：
  - `func TestXxx(t *testing.T)` — 測試函式命名規則
  - `t.Run("subtest", func(t *testing.T){})` — 子測試
  - `t.Errorf()` / `t.Fatalf()` — 報告錯誤
- **mockery 產生 mock** — 根據 interface 自動產生 mock struct，設定預期行為與回傳值
- **httptest** — 建立假的 HTTP 請求測試 handler，不需要真的啟動 server
- **Table-driven tests** — Go 社群慣用的測試風格，一個 slice 裝多組測試案例

#### ④ 動手做

**必做**
1. 執行 `go test ./...`，閱讀測試輸出
2. 在 `book_service_test.go` 新增一個測試方法：`TestFindByID_NotFound`

**延伸挑戰**
1. 在 `book_api_test.go` 新增測試：`TestUpdateBook_Success`
2. 用 `go test -cover ./...` 查看覆蓋率，嘗試提升至 60% 以上
3. 新增 404 與 Validation fail 測試案例

#### ⑤ 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| mock 方法沒被呼叫，測試仍通過 | 沒有用 `AssertExpectations` 驗證 | 在測試結尾加 `mockRepo.AssertExpectations(t)` |
| httptest 回傳空 body | ResponseRecorder 沒有正確讀取 | 用 `rec.Body.String()` 或 `rec.Body.Bytes()` 取得回應 |
| 測試間互相影響 | 共用了全域狀態 | 每個測試函式重新建立 mock 和待測物件 |

#### ⑥ 驗收標準（DoD）

能為新方法撰寫單元測試與 handler 測試，能解釋 mock 的用途。建議門檻：行覆蓋率 ≥ 60%，Service 核心流程覆蓋率 ≥ 80%。

---

## §6 里程碑檢核

| 里程碑 | 涵蓋單元 | 達成標準 | 未達標補救 |
|--------|----------|----------|-----------|
| **M1：能跑能看** | U01–U03 | 學員可啟動專案、做出假資料頁面與 API、在瀏覽器看到書籍清單 | 補做 U01–U03 的必做練習；講師一對一確認環境問題 |
| **M2：真實 CRUD 循環** | U04–U06 | 學員可串接真實 DB，頁面與 API 顯示 DB 資料，完成新增/查詢/更新/刪除 | 回頭確認 DB 連線設定，重做 U06 的 Repository 串接 |
| **M3：正規架構與完整功能** | U07–U10 | 學員有 Service/DTO 分層、統一錯誤處理、Redis 快取、自動化測試，可展示完整成品 | 優先完成 U07 分層重構；U09、U10 可簡化為閱讀理解 + 跑既有測試 |

---

## §7 每單元教學流程

### 標準流程（以 90 min 為基準）

| 時間 | 活動 | 負責 | 說明 |
|------|------|------|------|
| 10 min | 觀念講解 | 講師 | 核心概念 + 為什麼這樣設計（對應 ①③） |
| 20 min | 講師示範 | 講師 | 基於 Repo 程式碼 live demo |
| 40 min | 學員實作 | 學員為主、助教巡場 | 「必做」練習，助教協助個別卡關 |
| 15 min | Code Review | 講師主持 | 抽 2–3 人分享，討論不同寫法 |
| 5 min | 總結 | 講師 | 重點回顧 + 指派延伸挑戰作為課後作業 |

### 不同時數調整建議

| 時數 | 調整方式 |
|------|----------|
| 60 min（U02、U03、U04） | 壓縮觀念講解至 5 min、講師示範至 10 min、Code Review 至 5 min；學員實作維持 35 min |
| 120 min（U06、U10） | 多出的 30 min 分配至：講師示範 +10 min、學員實作 +20 min |

### 講師與助教分工

- **講師**：負責觀念講解、live demo、收斂共通問題、主持 Code Review
- **助教**：負責巡場協助個別卡關的學員、確認所有人環境正常、記錄高頻問題回饋給講師

### 進度異常處理

班級進度落後時，優先縮減的環節（建議順序）：
1. 總結（可移至下堂課開頭）
2. 講師示範（改為直接看 Repo 程式碼帶讀）
3. 觀念講解（精簡至核心要點）

**不可壓縮**：學員實作時間是底線——做中學是教案的核心設計。

---

## §8 班級分流建議

| 維度 | 基礎班 | 進階班 | 混合班 |
|------|--------|--------|--------|
| 練習範圍 | 只做必做 | 必做 + 全部延伸挑戰 | 課堂做必做，延伸作為課後作業 |
| 每單元時數 | 偏向 90–120 min | 可壓到 60–90 min | 90 min 為基準 |
| 自學比例 | 低（講師帶做） | 高（自行閱讀 + 實作） | 中等 |
| 可壓縮的單元 | 無（全部必修） | U02、U03 可合併為一堂快速帶過 | 視學員回饋彈性調整 |
| 課後追蹤 | 無延伸作業 | 延伸挑戰在下堂課開場 review | 延伸作業繳交 + 下堂課抽 review |

---

## §9 評量機制

### 課程練習評量

| 項目 | 比重 | 具體說明 |
|------|------|----------|
| 功能完成度 | 40% | 必做練習 100% 完成；能正確執行 CRUD、快取、驗證等核心功能 |
| 程式碼品質 | 20% | 命名規範（Go 慣例）、分層職責明確、無重複程式碼 |
| 測試覆蓋 | 20% | 具備正確的 Mock 運用與測試斷言；行覆蓋率 ≥ 60%，Service 核心流程 ≥ 80% |
| 問題分析 | 20% | 能判讀錯誤 log 並說明修復思路；遇到踩坑問題能自行排查 |

### 結訓驗收專案

> **設計理念**：驗收專案採用與教案**不同的業務主題**，驗證學員是否真正理解分層架構與開發流程，而非只是照抄書籍管理系統的程式碼。

**驗收方式**：限時實作（建議 3–4 小時）

**建議主題**：「員工通訊錄管理系統」（或其他具備 CRUD 性質的主題，如待辦事項、商品庫存）

| 驗收項目 | 達成條件 |
|----------|----------|
| Migration | 至少一支 up.sql + 一支 seed data SQL |
| Model + DTO | Model 對應 DB 表、CreateDTO 含 binding 驗證、ResponseDTO 含 ToResponse 轉換 |
| Repository | 完成基本 CRUD（GORM） |
| Service 層 | 介面與實作分離、業務邏輯集中在 Service |
| REST API | 至少 4 支 API（GET list / GET by ID / POST / DELETE），狀態碼正確 |
| Template 頁面 | 至少一個列表頁 + 一個新增表單頁，能正常操作 |
| 測試 | 至少 1 支 Service 單元測試 + 1 支 handler 測試 |
| 全部可運行 | `go run main.go` 能啟動、頁面能操作、API 能呼叫、`go test ./...` 全綠 |

**評分標準**：同上方課程練習評量比重。

### 標準化評語範例

| 等級 | 評語範例 |
|------|----------|
| 優秀 | 功能完整且程式碼結構清晰，測試覆蓋充分，能獨立排查問題並提出合理的設計理由 |
| 達標 | 核心功能完成，分層架構正確，有基本測試，能在提示下排查常見問題 |
| 待加強 | 部分功能缺失或分層不明確，測試不足，排查問題時需要較多協助 |
| 未達標 | 核心 CRUD 流程無法跑通，建議補做指定練習後重新驗收 |

---

## §10 參考文件

| 文件 | 用途 |
|------|------|
| `README.md` | 環境啟動與操作手冊 |
| `/swagger/index.html` | 互動式 API 測試介面 |

### 外部參考資源

| 主題 | 資源 |
|------|------|
| Go 官方文件 | [go.dev/doc](https://go.dev/doc/) |
| Effective Go | [go.dev/doc/effective_go](https://go.dev/doc/effective_go) |
| Gin 官方文件 | [gin-gonic.com/docs](https://gin-gonic.com/docs/) |
| GORM 官方文件 | [gorm.io/docs](https://gorm.io/docs/) |
| golang-migrate | [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) |
| html/template | [pkg.go.dev/html/template](https://pkg.go.dev/html/template) |
| go-redis | [redis.uptrace.dev](https://redis.uptrace.dev/) |
| Viper | [github.com/spf13/viper](https://github.com/spf13/viper) |
| swaggo/swag | [github.com/swaggo/swag](https://github.com/swaggo/swag) |
| Docker Compose 文件 | [docs.docker.com/compose](https://docs.docker.com/compose/) |

---

## §11 延伸學習方向

- **認證授權**：JWT middleware（golang-jwt）
- **進階查詢**：分頁、排序、模糊搜尋（GORM Scopes / Preload）
- **CI/CD**：GitHub Actions 自動測試與建置
- **雲端部署**：用 Docker 部署到 AWS / GCP / Azure
- **監控與可觀測性**：Prometheus metrics + Grafana、OpenTelemetry tracing
- **並發處理**：Go routines、channels、context 進階用法

---

## §12 維護與版本管理

### 教案版本資訊

| 項目 | 值 |
|------|-----|
| 教案版本 | v1.0 |
| 最後更新日期 | 2026-03-09 |
| 對應 Repo 分支 | `main` |
| Go 版本 | 1.22 |
| Gin 版本 | 1.10 |

### 與程式碼同步原則

- 當 Repo 的依賴版本升級（如 Go 版本、Gin 版本）時，教案需同步更新 §1 技術棧與 §4 驗證指令
- 當檔案路徑異動（重構、改名）時，教案 §5 各單元的「對應檔案」表格需同步更新
- 當 API 行為變更時，§5 U08 的動手做與踩坑提示需重新驗證
- 負責人：教案維護者應在每次 Repo 重大更新後檢查上述項目

### 開課前檢查清單

- [ ] §4 的驗證指令是否仍可正常執行？（`go version`、`docker compose version`、`go mod download`）
- [ ] `go test ./...` 是否全部通過？
- [ ] 所有單元列出的檔案路徑是否仍存在且正確？
- [ ] Swagger UI 是否可正常訪問？
- [ ] Docker Compose 是否能一鍵啟動所有服務？
- [ ] 驗收專案的模板或 starter 是否已準備好？
