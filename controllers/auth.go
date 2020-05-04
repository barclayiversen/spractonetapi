package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"spractonetapi/driver"
	"spractonetapi/models"
	"spractonetapi/repository/userRepository"
	"spractonetapi/utils"

	"strconv"

	"github.com/julienschmidt/httprouter"
)

func NewAuthController(tpl *template.Template) *Controller {
	return &Controller{tpl}
}

func (c Controller) VerifyEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(u)
	q := u.Query()
	fmt.Println(q["token"][0])
	fmt.Println(q["userid"][0])

	uuid := q["token"][0]
	var id int
	id, err = strconv.Atoi(q["userid"][0])
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}
	userRepo := userRepository.UserRepository{}
	err = userRepo.CheckKey(driver.DB, id, uuid)

	if err != nil {
		log.Println("Error creating email: ", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	var data models.Data
	data.Verified = true
	// http.Redirect(w, r, "/dashboard", 303)
	c.tpl.ExecuteTemplate(w, "emailverified.gohtml", data)
	return
}

//basic auth goes in here
func AlreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	return false
}
