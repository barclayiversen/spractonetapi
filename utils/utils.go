package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"spractonetapi/models"
	"strconv"
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

func NewMailSender(host string, port int, username, password string) *MailSender {
	return &MailSender{
		Addr: net.JoinHostPort(host, strconv.Itoa(port)),
		Auth: smtp.PlainAuth("", username, password, host),
		From: mail.Address{
			Name:    "Spracto net Mailer (no reply)",
			Address: "testenvsa2@gmail.com",
		},
	}
}

func (s *MailSender) Send(mail Mail) error {
	headers := map[string]string{
		"From":         s.From.String(),
		"To":           mail.To.String(),
		"Subject":      mail.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=utf-8",
	}
	msg := ""
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n"
	msg += mail.Body

	return smtp.SendMail(
		s.Addr,
		s.Auth,
		s.From.Address,
		[]string{mail.To.Address},
		[]byte(msg))
}

// Hoping to get this working later.
func SetHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in set header")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Authorization, Access-Control-Allow-Origin")
		next(w, r)
	}
}
