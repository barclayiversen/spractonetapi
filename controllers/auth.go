package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"spractonetapi/models"
	"spractonetapi/repository/userRepository"
	"spractonetapi/utils"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Controller struct{}

// TokenVerifyMiddleware ensures the integrity of the jwt.
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
				fmt.Println("token.valid: ", token.Valid)
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

func (c Controller) VerifyEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u, err := url.Parse(r.RequestURI)
		if err != nil {
			log.Println(err)
		}

		q := u.Query()
		uuid := q["token"][0]
		var id int
		id, err = strconv.Atoi(q["userid"][0])
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid id")
			return
		}

		userRepo := userRepository.UserRepository{}
		err = userRepo.CheckKey(db, id, uuid)
		if err != nil {
			log.Println("Error creating email: ", err)
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		var user models.User
		user, err = userRepo.GetUserById(db, user, id)
		token, err := utils.GenerateToken(user)
		if err != nil {
			utils.ResponseJSON(w, "Email Verified, error getting token. Please log in.")
		}

		utils.ResponseJSON(w, token)
		return
	}
}
