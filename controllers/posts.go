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
	"github.com/gorilla/mux"
)

// GetUserPosts gets posts by user ID
func (c Controller) GetUserPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
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
		id := int(claims["sub"].(float64))
		postRepo := postRepository.PostRepository{}
		posts, err := postRepo.GetUserPosts(db, id)
		if err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusBadRequest, "Error getting user's posts")
		}
		utils.ResponseJSON(w, posts)
	}
}

// CreatePost creates a new post by the user identified in the token
func (c Controller) CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
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

		var post models.Post
		json.NewDecoder(r.Body).Decode(&post)

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
		post, err := postRepo.CreatePost(db, post)
		if err != nil {
			fmt.Println("Create post error:", err)
			utils.RespondWithError(w, http.StatusBadRequest, "Error creating post")
		}

		utils.ResponseJSON(w, post)
		return
	}
}

// DeletePost deletes a post based on the post id in the params and the user id in the token.
func (c Controller) DeletePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromToken(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "This was very unexpected")
			return
		}
		params := mux.Vars(r)
		postID := params["id"]

		postRepo := postRepository.PostRepository{}
		err = postRepo.DeletePost(db, userID, postID)
		if err != nil {
			fmt.Println(err)
			return
		}
		utils.ResponseJSON(w, "Post Deleted")
	}
}
