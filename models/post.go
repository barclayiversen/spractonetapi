package models

type Post struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Post       string `json:"post"`
	Created_at int    `json:"created_at"`
	User_id    int    `json:"user_id"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}
