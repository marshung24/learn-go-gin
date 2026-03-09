-- Seed data：預先塞入幾筆範例書籍，方便開發與 demo 時直接有資料可看。
-- 正式環境可選擇不執行此 migration。

INSERT INTO books (title, author, isbn, stock) VALUES
('Clean Code',              'Robert C. Martin', '9780132350884', 5),
('Effective Java',          'Joshua Bloch',     '9780134685991', 3),
('Go Programming Language', 'Alan Donovan',     '9780134190440', 8);
