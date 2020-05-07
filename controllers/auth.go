package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"spractonetapi/repository/userRepository"
	"spractonetapi/utils"
	"strconv"
)

func (c Controller) VerifyEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		err = userRepo.CheckKey(db, id, uuid)

		if err != nil {
			log.Println("Error creating email: ", err)
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseJSON(w, "Email Verified!")
		return
	}
}
