package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"spractonetapi/models"
	"spractonetapi/repository/postRepository"
	"spractonetapi/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// GetUserPosts gets posts by user ID
func (c Controller) GetUserPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, "End point works")
	}
}

// CreatePost creates a new post by the user identified in the token
func (c Controller) CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		fmt.Println("can I access headers in here? ", authHeader)

		bearerToken := strings.Split(authHeader, " ")
		authToken := bearerToken[1]

		token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if error != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, error.Error())
			return
		}
		claims := token.Claims.(jwt.MapClaims)

		fmt.Printf("%T\n", claims["sub"])
		fmt.Println(claims["sub"])

		var post models.Post
		json.NewDecoder(r.Body).Decode(&post)
		fmt.Println(post)
		if post.Title == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Posts require a title")
			return
		}

		if post.Post == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Posts cannot be blank!")
			return
		}

		post.User_id = int(claims["sub"].(float64))
		postRepo := postRepository.PostRepository{}
		err := postRepo.CreatePost(db, post)
		if err != nil {
			fmt.Println("Create post error:", err)
			utils.RespondWithError(w, http.StatusBadRequest, "Error creating post")
		}

		utils.ResponseJSON(w, "Post Created")
		return
	}
}
