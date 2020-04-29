package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"spractonetapi/controllers"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	godotenv.Load()
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

}

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(tpl)

	r.GET("/", uc.IndexHandler)
	r.GET("/signup", uc.Signup)
	r.POST("/signup", uc.Signup)
	r.GET("/login", uc.Login)
	r.POST("/login", uc.Login)

	r.ServeFiles("/static/*filepath", http.Dir("./static"))
	fmt.Println("listening on 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
