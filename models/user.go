package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Username string `json:"username"`
}

type Hobby struct {
	Hobby string `json:"hobby"`
}
