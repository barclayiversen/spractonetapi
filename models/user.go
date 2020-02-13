package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Token    string `json:"token"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Username string `json:"username" validate:"required"`
}
