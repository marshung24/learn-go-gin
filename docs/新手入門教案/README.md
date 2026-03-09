# Go 新手入門教案

> 以本 Repo「書籍管理系統」為實作主軸，從零帶出 Go + Gin 全棧開發核心技能。

---

## §1 教案資訊

| 項目 | 說明 |
|------|------|
| 對象 | 有基礎程式概念（變數、迴圈、函式）與基本命令列操作經驗，但**尚未做過 Go Web 專案**的初學者 |
| 目標 | 獨立完成一個「能跑、能測、能部署」的 CRUD Web 應用（書籍管理系統） |
| 技術棧 | Go 1.22、Gin 1.10、GORM、golang-migrate、MySQL 8、Redis 7、html/template、swaggo/swag（Swagger）、Docker Compose |
| 時數 | 課前準備（自學）+ 10 單元（各 60–120 分鐘） |
| 教學方式 | 觀念講解 → Repo 實作示範 → 學員動手做 → Code Review → 驗收 |
| 先備知識 | 基礎 SQL（SELECT / INSERT）、HTTP 概念（GET / POST）、Git 基本操作（clone / commit） |
| 不含範圍 | 前端框架（React / Vue）、JWT 認證、gRPC、微服務架構、Kubernetes |

---

## §2 學習成果（結訓時能做到）

1. **能說明** Gin 專案分層架構（Router/Handler → Service → Repository → Model）各層的職責與呼叫方向
2. **能使用** Go Modules 建置專案、用 Docker Compose 一鍵啟動整套開發環境（App + MySQL + Redis）
3. **能撰寫** golang-migrate migration 腳本管理資料庫 schema 演進
4. **能實作** GORM Repository 完成 CRUD 操作，新增自訂查詢方法
5. **能設計** RESTful API（正確使用 HTTP 動詞與狀態碼），用 middleware 統一錯誤回應，並產出 Swagger 互動式文件
6. **能製作** html/template 伺服器端渲染頁面（列表、表單、詳情），理解 PRG 模式
7. **能實作** Redis Cache-Aside 快取模式，並用 redis-cli 驗證快取行為
8. **能撰寫** Go testing 單元測試與 handler 測試，使用 mock interface 驗證依賴替身

---

## §3 閱讀指引

### 練習標記

| 標記 | 意義 |
|------|------|
| **必做** | 課堂內所有學員都應完成的核心練習 |
| **延伸挑戰** | 給進度快的學員或課後作業 |

### 依賴標記

- `⟵ 需完成 UXX` = 需先完成指定單元的對應練習
- `⟵ 需完成 UXX 延伸` = 特指依賴某單元的「延伸挑戰」成果

### 命名規範

| 項目 | 格式 | 範例 |
|------|------|------|
| 單元編號 | U00–U10 | U03 |
| 檔案路徑 | 省略 `internal/` | `handler/book_api.go` |
| Migration 路徑 | 完整路徑 | `migrations/000001_create_books_table.up.sql` |
| 測試路徑 | 省略 `internal/` | `service/book_service_test.go` |

### 單元內部結構

每個單元依序包含六個子區塊：① 為什麼先教這個？ → ② 對應檔案 → ③ 核心觀念 → ④ 動手做 → ⑤ 踩坑提示 → ⑥ 驗收標準（DoD）

---

## 課程目錄

| 單元 | 名稱 | 時數 | 一句話目標 | 前置依賴 |
|------|------|------|-----------|----------|
| [U00](00-課前準備.md) | 課前準備 | 自學 | 安裝開發環境並驗證可用 | — |
| [U01](01-專案啟動與第一支API.md) | 專案啟動與第一支 API | 90 min | 能啟動專案並用瀏覽器呼叫自訂端點 | — |
| [U02](02-假畫面與Template初體驗.md) | 假畫面 — Template 初體驗 | 60 min | 能用 html/template 做出書籍清單頁面（假資料） | U01 |
| [U03](03-API回傳假資料.md) | API 回傳假資料 | 60 min | 能用 Gin handler 回傳固定 JSON 並在 Swagger 測試 | U01 |
| [U04](04-設定檔與多環境切換.md) | 設定檔與多環境切換 | 60 min | 能解釋環境變數機制並用 Docker Compose 啟動完整環境 | U01 |
| [U05](05-golang-migrate與資料庫設計.md) | golang-migrate 與資料庫設計 | 90 min | 能撰寫 migration SQL 並執行 up/down | U04 |
| [U06](06-GORM串接資料庫.md) | GORM 串接資料庫 | 120 min | 能用 GORM 完成真實 CRUD，頁面與 API 顯示 DB 真實資料 | U02, U03, U05 |
| [U07](07-正規化架構與Service分層.md) | 正規化架構 — Service 與 DTO 分層 | 90 min | 能從 Handler 抽出 Service 層，引入 DTO 做資料驗證與格式轉換 | U06 |
| [U08](08-REST-API完善與錯誤處理.md) | REST API 完善與錯誤處理 | 90 min | 能用 middleware 統一錯誤回應，正確使用 HTTP 狀態碼 | U07 |
| [U09](09-Redis快取整合.md) | Redis 快取整合 | 90 min | 能用 redis-cli 驗證快取行為並說明 Cache-Aside 流程 | U07 |
| [U10](10-測試策略與實戰.md) | 測試策略與實戰 | 120 min | 能為新方法撰寫單元測試與 handler 測試 | U08 |

> **設計說明**：單元順序遵循「對應專案流程的漸進式推演」——環境確認（U01）→ 假畫面 Level-1 UI（U02）→ 假資料 Level-2 程式（U03）→ 設定環境（U04）→ 資料庫設計 Level-3 資料（U05）→ 串接真實 DB（U06）→ 架構正規化（U07）→ 完善 API（U08）→ 效能優化（U09）→ 測試驗證（U10）。

---

## §6 里程碑檢核

| 里程碑 | 涵蓋單元 | 達成標準 | 未達標補救 |
|--------|----------|----------|-----------|
| **M1：能跑能看** | U01–U03 | 學員可啟動專案、做出假資料頁面與 API、在瀏覽器看到書籍清單 | 補做 U01–U03 的必做練習；講師一對一確認環境問題 |
| **M2：真實 CRUD 循環** | U04–U06 | 學員可串接真實 DB，頁面與 API 顯示 DB 資料，完成新增/查詢/更新/刪除 | 回頭確認 DB 連線設定，重做 U06 的 Repository 串接 |
| **M3：正規架構與完整功能** | U07–U10 | 學員有 Service/DTO 分層、統一錯誤處理、Redis 快取、自動化測試 | 優先完成 U07 分層重構；U09、U10 可簡化為閱讀理解 + 跑既有測試 |

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

- **講師**：觀念講解、live demo、收斂共通問題、主持 Code Review
- **助教**：巡場協助個別卡關的學員、確認所有人環境正常、記錄高頻問題回饋給講師

### 進度異常處理

班級進度落後時，優先縮減（依順序）：總結 → 講師示範 → 觀念講解。**不可壓縮**：學員實作時間是底線。

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
| 程式碼品質 | 20% | 命名規範、分層職責明確、無重複程式碼 |
| 測試覆蓋 | 20% | 具備正確的 Mock 運用與測試斷言；行覆蓋率 ≥ 60% |
| 問題分析 | 20% | 能判讀錯誤 log 並說明修復思路 |

### 結訓驗收專案

> 驗收專案採用與教案**不同的業務主題**（如「員工通訊錄管理系統」），驗證學員是否真正理解開發流程。

**驗收方式**：限時實作（建議 3–4 小時）

| 驗收項目 | 達成條件 |
|----------|----------|
| Migration | 至少一支建表 SQL + 一支 seed data SQL |
| Model + DTO | Model 對應 DB 表、CreateDTO 含 binding 驗證、ResponseDTO 含轉換方法 |
| Repository | 完成基本 CRUD（GORM） |
| Service 層 | 介面與實作分離、業務邏輯集中在 Service |
| REST API | 至少 4 支 API（GET list / GET by ID / POST / DELETE），狀態碼正確 |
| Template 頁面 | 至少一個列表頁 + 一個新增表單頁 |
| 測試 | 至少 1 支 Service 單元測試 + 1 支 handler 測試 |
| 全部可運行 | `go run main.go` 能啟動、頁面能操作、API 能呼叫、`go test ./...` 全綠 |

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
| Gin 官方文件 | [gin-gonic.com/docs](https://gin-gonic.com/docs/) |
| GORM 官方文件 | [gorm.io/docs](https://gorm.io/docs/) |
| golang-migrate | [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) |
| html/template | [pkg.go.dev/html/template](https://pkg.go.dev/html/template) |
| Docker Compose 文件 | [docs.docker.com/compose](https://docs.docker.com/compose/) |

---

## §11 延伸學習方向

- **認證授權**：JWT middleware（如 golang-jwt）
- **進階查詢**：分頁、排序、模糊搜尋（GORM Scopes）
- **CI/CD**：GitHub Actions 自動測試與建置
- **雲端部署**：Docker 部署到 AWS / GCP / Azure
- **監控**：Prometheus metrics + Grafana

---

## §12 維護與版本管理

| 項目 | 值 |
|------|-----|
| 教案版本 | v1.0 |
| 最後更新日期 | 2026-03-09 |
| 對應 Repo 分支 | `main` |
| Go 版本 | 1.22 |
| Gin 版本 | 1.10 |

### 開課前檢查清單

- [ ] §4 的驗證指令是否仍可正常執行？
- [ ] `go test ./...` 是否全部通過？
- [ ] 所有單元列出的檔案路徑是否仍存在且正確？
- [ ] Swagger UI 是否可正常訪問？
- [ ] Docker Compose 是否能一鍵啟動所有服務？
