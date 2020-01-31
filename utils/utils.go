package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"spractonetapi/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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
