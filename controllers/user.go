package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"spractonetapi/models"
	"spractonetapi/repository/userRepository"
	"spractonetapi/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Authorization, Access-Control-Allow-Origin")

		if r.Method == "OPTIONS" {
			return
		}

		var user models.User
		json.NewDecoder(r.Body).Decode(&user)
		fmt.Println("user in signup ", user)

		if user.Email == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Email is missing")
			return
		}

		if user.Password == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Password is missing")
			return
		}

		if user.Username == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Username is missing")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)

		userRepo := userRepository.UserRepository{}
		user, err = userRepo.Signup(db, user)

		//This is bad
		if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			utils.RespondWithError(w, http.StatusBadRequest, "That email is already registered")
			return
		}

		if err != nil && err.Error() != "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		}

		token, err := utils.GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		//Data structures could be better here
		user.Password = ""
		user.Token = token
		utils.ResponseJSON(w, user)
	}

}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") //find a way to auto add this header in every request.
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Authorization, Access-Control-Allow-Origin")

		if r.Method == "OPTIONS" {
			return
		}

		var user models.User

		json.NewDecoder(r.Body).Decode(&user)

		password := user.Password

		userRepo := userRepository.UserRepository{}
		user, err := userRepo.Login(db, user)

		hashedPassword := user.Password

		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondWithError(w, http.StatusBadRequest, "The User does not exist")
				return
			} else {
				log.Fatal(err)
			}
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Password")
			return
		}

		token, err := utils.GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		user.Token = token
		user.Password = ""
		utils.ResponseJSON(w, user)

	}
}

func (c Controller) TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Headers", "Accepts, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			return
		}

		var errorObject models.Error
		authHeader := r.Header.Get("Authorization")

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {

			authToken := bearerToken[1]
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, error.Error())
				return
			}

			if token.Valid {

				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, error.Error())
				return
			}
		} else {
			errorObject.Message = "Invalid Token"
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}
	})
}
