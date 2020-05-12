package postRepository

import (
	"database/sql"
	"fmt"
	"spractonetapi/models"
	"strconv"
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

// GetUserPosts gets a user's posts by userid
func (p PostRepository) GetUserPosts(db *sql.DB, userID int) ([]models.Post, error) {
	var post models.Post
	P := []models.Post{}
	fmt.Println(userID)
	stmt := "SELECT * FROM posts WHERE user_id = $1"
	rows, err := db.Query(stmt, userID)
	if err != nil {
		fmt.Println(err)
		return P, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Title, &post.Post, &post.User_id, &post.Created_at)
		if err != nil {
			fmt.Println(err)
			return P, err
		}
		P = append(P, post)
	}

	fmt.Println(P)
	return P, nil
}

func (p PostRepository) DeletePost(db *sql.DB, userId int, postId string) error {
	stmt := "SELECT user_id FROM posts WHERE id = $1"
	id, err := strconv.Atoi(postId)
	if err != nil {
		return &MyError{}
	}
	row := db.QueryRow(stmt, id)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	var post models.Post
	row.Scan(&post.User_id)

	if userId != post.User_id {
		return &MyError{}
	}
	fmt.Println(row)
	return nil
}

type MyError struct{}

func (myErr *MyError) Error() string {
	return "Something unexpected happened"
}
