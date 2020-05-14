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
func (p PostRepository) CreatePost(db *sql.DB, post models.Post) (models.Post, error) {
	post.Created_at = int(time.Now().Unix())
	stmt := "INSERT INTO posts (title, post, created_at, user_id) VALUES ($1, $2, $3, $4) RETURNING created_at"
	_, err := db.Exec(stmt, post.Title, post.Post, post.Created_at, post.User_id)
	if err != nil {
		fmt.Println(err)
		return post, err
	}

	return post, nil
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
	stmt := "SELECT * FROM posts WHERE id = $1"

	id, err := strconv.Atoi(postId)
	if err != nil {
		fmt.Println("err 1", err)
		//return &MyError{}
		return err
	}
	fmt.Println("thisx is a s")
	row := db.QueryRow(stmt, id)
	var post models.Post
	row.Scan(&post.Id, &post.Title, &post.Post, &post.User_id, &post.Created_at)

	if userId != post.User_id {
		fmt.Println("err 2", err)
		return &MyError{}
		//return err
	}
	stmt = "DELETE FROM posts WHERE id = $1"
	_, err = db.Exec(stmt, postId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type MyError struct{}

func (myErr *MyError) Error() string {
	return "Something unexpected happened"
}
