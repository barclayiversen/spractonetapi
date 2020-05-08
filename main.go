package main

import (
	"database/sql"
	"log"
	"net/http"

	"spractonetapi/controllers"
	"spractonetapi/driver"

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
	router.HandleFunc("/", controller.SetHeader(controller.HelloWorld())).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", controller.SetHeader(controller.Login(db))).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/{id}", controller.SetHeader(controller.TokenVerifyMiddleware(controller.GetUserById(db)))).Methods("GET", "OPTIONS")
	router.HandleFunc("/users", controller.SetHeader(controller.Signup(db))).Methods("POST", "OPTIONS")
	router.HandleFunc("/verifyemail", controller.SetHeader(controller.VerifyEmail(db))).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}/posts", controller.SetHeader(controller.TokenVerifyMiddleware(controller.GetUserPosts(db)))).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}/posts", controller.SetHeader(controller.TokenVerifyMiddleware(controller.CreatePost(db)))).Methods("POST", "OPTIONS")

	router.Use(mux.CORSMethodMiddleware(router))

	log.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
