package models

type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	UserName   string `json:"user_name"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfileURL string `json:"profile_url"`
	City       string `json:"city"`
}
