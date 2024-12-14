package models

// User structure
type User struct {
	ID    int64  `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
