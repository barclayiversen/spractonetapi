package models

type Post struct {
	Title      string `json:"title"`
	Post       string `json:"post"`
	Created_at int    `json:"created_at"`
	User_id    int    `json:"user_id"`
}
