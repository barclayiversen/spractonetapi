package main

import (
	"database/sql"
	"log"
	"net/http"

	"goapi/controllers"
	"goapi/driver"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var db *sql.DB

func init() {
	godotenv.Load()
}

func main() {

	db = driver.ConnectDB()
	controller := controllers.Controller{}
	router := mux.NewRouter()

	router.HandleFunc("/test", controller.TokenVerifyMiddleware(controller.Test()))
	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/{id}", controller.TokenVerifyMiddleware(controller.GetUserById(db))).Methods("GET", "OPTIONS")
	// router.HandleFunc("/users/{id}", controller.TokenVerifyMiddleware(controller.SaveUserData(db))).Methods("POST", "OPTIONS")
	router.HandleFunc("/protected", controller.TokenVerifyMiddleware(controller.ProtectedEndpoint()))

	router.Use(mux.CORSMethodMiddleware(router))

	log.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
