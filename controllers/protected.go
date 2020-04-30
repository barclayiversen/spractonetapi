package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"spractonetapi/models"
	"spractonetapi/repository/userRepository"
	"spractonetapi/utils"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/gorilla/mux"
)

//type Controller struct{}

func NewProtectedController(tpl *template.Template) *Controller {
	return &Controller{tpl}
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

func (c Controller) DashHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// check for existence of cookie

	// if cookie exists send it to a local redis instance and see if it exists there.
	// if it does, the value to that key is the user id, possibly other info that
	// would be kept server side.
	c.tpl.ExecuteTemplate(w, "dashboard.gohtml", nil)
}
