package dto

// BookCreateDTO 書籍新增/更新請求 DTO
// binding tag 用於 Gin 的請求驗證
type BookCreateDTO struct {
	Title  string `json:"title" binding:"required,min=1,max=200"`
	Author string `json:"author" binding:"required,min=1,max=100"`
	ISBN   string `json:"isbn" binding:"required,max=20"`
	Stock  int    `json:"stock" binding:"gte=0"`
}
