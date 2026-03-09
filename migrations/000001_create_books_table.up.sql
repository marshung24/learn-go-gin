-- ==========================================================
-- golang-migrate 命名慣例：{序號}_{描述}.{up|down}.sql
--   序號遞增（000001, 000002...），up 是正向，down 是回滾。
--   每個 up.sql 都要有對應的 down.sql，確保能回滾。
-- ==========================================================
-- 欄位設計說明：
--   id         — 自增主鍵，BIGINT 避免溢位
--   isbn       — UNIQUE 確保同一本書不會重複匯入
--   stock      — 預設 0，代表「尚未進貨」
--   created_at / updated_at — 自動記錄時間戳，updated_at 在每次 UPDATE 自動更新
-- ==========================================================

CREATE TABLE books (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    title      VARCHAR(200)    NOT NULL,
    author     VARCHAR(100)    NOT NULL,
    isbn       VARCHAR(20)     NOT NULL UNIQUE,
    stock      INT             NOT NULL DEFAULT 0,
    created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
