package postRepository

import (
	"database/sql"
	"fmt"
	"spractonetapi/models"
	"time"
)

type PostRepository struct{}

// CreatePost is a function
func (p PostRepository) CreatePost(db *sql.DB, post models.Post) error {
	createdAt := time.Now().Unix()
	stmt := "INSERT INTO posts (title, post, created_at, user_id) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(stmt, post.Title, post.Post, createdAt, post.User_id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
