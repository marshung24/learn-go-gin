# U05｜golang-migrate 與資料庫設計

> 能撰寫 migration SQL 並執行 up/down ｜ 90 min ｜ 前置依賴：U04

---

## 為什麼先教這個？

U02-U03 用的是假資料，接下來要換成真的。第一步是「設計資料庫結構」——用 golang-migrate 管理 schema 版本，確保多人協作時每個人的 DB 長一樣。

---

## 對應檔案

| 檔案 | 角色 |
|------|------|
| `migrations/000001_create_books_table.up.sql` | 建表（初始 schema） |
| `migrations/000001_create_books_table.down.sql` | 回滾建表 |
| `migrations/000002_insert_sample_data.up.sql` | 塞入範例資料（seed data） |
| `migrations/000002_insert_sample_data.down.sql` | 回滾範例資料 |
| `Makefile` | migration 指令封裝 |

---

## 核心觀念

### 命名慣例

```
{序號}_{描述}.{up|down}.sql
```

- 序號遞增（000001, 000002, ...）
- up 是正向 migration
- down 是回滾

### `schema_migrations` 表

golang-migrate 在 DB 中自動建立的紀錄表，記錄當前版本：

```sql
SELECT * FROM schema_migrations;
-- version | dirty
-- 2       | false
```

### up/down 對稱

每個 up.sql 都要有對應的 down.sql，確保能回滾：

```sql
-- 000001_create_books_table.up.sql
CREATE TABLE books (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    author VARCHAR(100) NOT NULL,
    isbn VARCHAR(20) UNIQUE,
    stock INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 000001_create_books_table.down.sql
DROP TABLE IF EXISTS books;
```

### 欄位設計思考

| 設計決策 | 理由 |
|----------|------|
| `BIGINT` 而非 `INT` | 預留空間，避免未來 ID 溢出 |
| `isbn UNIQUE` | ISBN 唯一性約束，防止重複 |
| `DEFAULT CURRENT_TIMESTAMP` | 建立時自動填入時間 |
| `ON UPDATE CURRENT_TIMESTAMP` | 更新時自動更新時間 |

### `migrate` CLI

```bash
# 執行所有 up migration
migrate -path migrations -database "mysql://app:app@tcp(localhost:3306)/app" up

# 回滾一版
migrate -path migrations -database "mysql://..." down 1

# 查看當前版本
migrate -path migrations -database "mysql://..." version
```

或用 Makefile 封裝：

```bash
make migrate-up
make migrate-down
```

---

## 動手做

### 必做

**1. 閱讀現有 migration**

閱讀 000001 和 000002 的 SQL，說出每個欄位的設計理由。

**2. 新增 migration**

新增一組 migration，為 books 表加一個 `publisher VARCHAR(100)` 欄位：

```sql
-- 000003_add_publisher_column.up.sql
ALTER TABLE books ADD COLUMN publisher VARCHAR(100);

-- 000003_add_publisher_column.down.sql
ALTER TABLE books DROP COLUMN publisher;
```

**3. 執行並驗證**

```bash
make migrate-up
```

連線 MySQL 查詢 `SELECT * FROM schema_migrations` 確認版本。

### 延伸挑戰

1. 執行 `make migrate-down` 回滾一版，確認 publisher 欄位消失；思考：什麼情況下會需要 down migration？

---

## 踩坑提示

| 現象 | 原因 | 解法 |
|------|------|------|
| `dirty database version` | 上次 migration 執行失敗但版本號已更新 | 手動修正 `schema_migrations` 表的 dirty 欄位為 false，或 `migrate force VERSION` |
| migration 檔案被忽略 | 序號不連續或格式錯誤 | 確認序號遞增且格式為 `{序號}_{描述}.{up\|down}.sql` |

---

## 驗收標準（DoD）

- [ ] 能撰寫一組新的 migration SQL（up + down）並成功執行
- [ ] 能在 `schema_migrations` 確認版本紀錄
- [ ] 能說出 up 和 down migration 的關係
