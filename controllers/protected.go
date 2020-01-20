package controllers

import (
	"database/sql"
	"fmt"
	"goapi/models"
	"goapi/repository/userRepository"
	"goapi/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Controller struct{}

func (c Controller) ProtectedEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("\"yes\""))
		fmt.Println("protected endpoint invoked")
	}
}

func (c Controller) Test() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, "testing")
	}
}

func (c Controller) GetUserById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Access-Control-Allow-Origin")
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
