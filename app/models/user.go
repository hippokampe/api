package models

type User struct {
	ID       string `json:"id"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Username string `json:"username"`
}
