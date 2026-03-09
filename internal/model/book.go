package model

import (
	"time"
)

// Book 對應資料庫 books 資料表
// GORM 透過 tag 對應欄位名稱與約束
type Book struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"column:title;size:200;not null" json:"title"`
	Author    string    `gorm:"column:author;size:100;not null" json:"author"`
	ISBN      string    `gorm:"column:isbn;size:20;uniqueIndex;not null" json:"isbn"`
	Stock     int       `gorm:"column:stock;default:0;not null" json:"stock"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName 指定資料表名稱
// GORM 預設會將 struct 名稱轉為複數形式，這裡明確指定
func (Book) TableName() string {
	return "books"
}
