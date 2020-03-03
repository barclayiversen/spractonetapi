package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"spractonetapi/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/validator.v9"
)

type MailSender struct {
	Addr string
	Auth smtp.Auth
	From mail.Address
}

// Mail to send.
type Mail struct {
	To      mail.Address
	Subject string
	Body    string
}

func IsEmail(u models.User) bool {
	v := validator.New()
	err := v.Struct(u)

	if err != nil {
		return false
	}

	return true
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	var error models.Error
	error.Message = message
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func GenerateToken(user models.User) (string, error) {

	var err error
	secret := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(60 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"iss":      "spractonet",
		"exp":      expirationTime,
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func Send(u models.User, activationUrl string) error {
	auth := smtp.PlainAuth("", os.Getenv("FROM"), os.Getenv("MAILPASS"), os.Getenv("SMTPSERVER"))

	to := []string{u.Email}
	msg := []byte("From: " + "\n" + "To:" + u.Email + "\r\n" +
		"Subject: Welcome to Spracto net \r\n" +
		"\r\n" +
		"Here's your link to activate your account: " + activationUrl)

	err := smtp.SendMail(os.Getenv("SMTPSERVER")+":587", auth, "noreply@spracto.net", to, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}
	log.Print("sent " + activationUrl)

	return nil
}
