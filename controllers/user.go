package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"spractonetapi/driver"
	"spractonetapi/models"
	"spractonetapi/repository/userRepository"
	"spractonetapi/utils"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	tpl *template.Template
}

func NewUserController(tpl *template.Template) *Controller {
	return &Controller{tpl}
}

func (c Controller) HelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func (c Controller) IndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	c.tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func (c Controller) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == "GET" {
		c.tpl.ExecuteTemplate(w, "login.gohtml", nil)
		return
	}

	if r.Method == "POST" {
		var user models.User

		json.NewDecoder(r.Body).Decode(&user)

		password := user.Password

		userRepo := userRepository.UserRepository{}
		user, err := userRepo.Login(driver.DB, user)

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

func (c Controller) Signup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if r.Method == "GET" {
		c.tpl.ExecuteTemplate(w, "signup.gohtml", nil)
		return
	}

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
	user, err = userRepo.Signup(driver.DB, user)

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

func (c Controller) TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Headers", "Accepts, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")

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
