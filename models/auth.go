package models

type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	ID         string
	Email      string
	UserName   string
	FirstName  string
	LastName   string
	ProfileURL string
	City       string
}
