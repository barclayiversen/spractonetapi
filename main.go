package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"spractonetapi/controllers"
	"spractonetapi/driver"

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
	ac := controllers.NewAuthController(tpl)
	p := controllers.NewProtectedController(tpl)
	driver.ConnectDB()
	r.GET("/", uc.IndexHandler)
	r.GET("/dashboard", p.DashHandler)
	r.GET("/signup", uc.Signup)
	r.POST("/signup", uc.Signup)
	r.GET("/login", uc.Login)
	r.POST("/login", uc.Login)
	r.GET("/verifyemail", ac.VerifyEmail)

	r.ServeFiles("/static/*filepath", http.Dir("./static"))
	fmt.Println("listening on 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
