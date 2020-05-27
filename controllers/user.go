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
	"strconv"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (c Controller) HelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, "Welcome to Spracto net!")
	}
}

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User
		json.NewDecoder(r.Body).Decode(&user)

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

		EmailIsValid := utils.IsEmail(user)

		if EmailIsValid != true {
			utils.RespondWithError(w, http.StatusBadRequest, "Please provide a valid Email address IE email@domain.com")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong. Please try again later.")
			return
		}

		user.Password = string(hash)

		u2 := uuid.Must(uuid.NewV4()).String()
		if err != nil {
			fmt.Printf("Something went wrong: %s", err)
			return
		}

		user.SignupKey = u2
		userRepo := userRepository.UserRepository{}
		user, err = userRepo.Signup(db, user)

		//This is bad
		if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			utils.RespondWithError(w, http.StatusBadRequest, "That email is already registered")
			return
		}

		if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"username_is_unique\"" {
			utils.RespondWithError(w, http.StatusBadRequest, "That username is taken")
			return
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		// Create magic link
		var tokenString string

		if os.Getenv("ENV") == "dev" {
			tokenString = "http://localhost:8080/verifyemail?token="
		}

		if os.Getenv("ENV") == "prod" {
			tokenString = "https://api.spracto.net/verifyemail?token="
		}

		fmt.Printf("UUIDv4: %s\n", u2)
		tokenString += u2
		tokenString += "&userid="
		tokenString += strconv.Itoa(user.ID)

		fmt.Println(tokenString)

		err = utils.Send(user, tokenString)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadGateway, "We weren't able to send you a verifcation email.")
			return
		}

		return
	}

}

func (c Controller) Login(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

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
				log.Println(err)
				utils.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
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

// GetUserById is for the initial population of data in the user dashboard
func (c Controller) GetUserById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Access-Control-Allow-Origin")
		var user models.User
		params := mux.Vars(r)
		userRepo := userRepository.UserRepository{}
		userId, err := strconv.Atoi(params["id"])
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}
		user, err = userRepo.GetUserById(db, user, userId)
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}
		utils.ResponseJSON(w, user)
	}
}

func (c Controller) UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, "Update User works")
	}
}

func (c Controller) DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, "Delete User works")
	}
}
