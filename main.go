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
	router.HandleFunc("/verifyemail", controller.SetHeader(controller.VerifyEmail(db))).Methods("GET", "OPTIONS")

	// Session source
	router.HandleFunc("/login", controller.SetHeader(controller.Login(db))).Methods("POST", "OPTIONS")
	// Users resource
	router.HandleFunc("/users", controller.SetHeader(controller.Signup(db))).Methods("POST", "OPTIONS")                                                    // Create
	router.HandleFunc("/users/{id}", controller.SetHeader(controller.TokenVerifyMiddleware(controller.GetUserById(db)))).Methods("GET", "OPTIONS")         // Read
	router.HandleFunc("/users/{id}", controller.SetHeader(controller.TokenVerifyMiddleware(controller.UpdateUser(db)))).Methods("PUT", "PATCH", "OPTIONS") // Update
	router.HandleFunc("/users/{id}", controller.SetHeader(controller.TokenVerifyMiddleware(controller.DeleteUser(db)))).Methods("DELETE", "OPTIONS")       // Destroy

	// Posts resource
	router.HandleFunc("/users/{id}/posts", controller.SetHeader(controller.TokenVerifyMiddleware(controller.GetUserPosts(db)))).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}/posts", controller.SetHeader(controller.TokenVerifyMiddleware(controller.CreatePost(db)))).Methods("POST", "OPTIONS")
	router.HandleFunc("/posts/{id}", controller.SetHeader(controller.TokenVerifyMiddleware(controller.DeletePost(db)))).Methods("DELETE", "OPTIONS")

	router.Use(mux.CORSMethodMiddleware(router))
	log.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
