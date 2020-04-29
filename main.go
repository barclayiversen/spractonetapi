package main

import (
	"context"
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

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, "params", ps)

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "login.gohtml", nil)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	}
}
