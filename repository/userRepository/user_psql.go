package userRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"spractonetapi/models"
	"time"
)

type UserRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Signup commits user data to the database and returns a user object
func (u UserRepository) Signup(db *sql.DB, user models.User) (models.User, error) {
	stmt := "INSERT INTO USERS (email,password,username,activation_key, created_at, activated) VALUES ($1, $2, $3, $4, $5, false) RETURNING id;"
	createdAt := time.Now().Unix()
	err := db.QueryRow(stmt, user.Email, user.Password, user.Username, user.SignupKey, createdAt).Scan(&user.ID)

	if err != nil {
		fmt.Println(err)
		return user, err
	}

	user.Password = ""
	return user, nil
}

// Login returs user data to be used for authentication.
func (u UserRepository) Login(db *sql.DB, user models.User) (models.User, error) {

	row := db.QueryRow("SELECT id, email, username, password, activated FROM users WHERE email = $1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Activated)

	if err != nil {
		return user, err
	}
	if user.Activated == false {
		return user, errors.New("Please activate your account to log in.")
	}

	return user, nil
}

// GetUserById retrieves the specified user
func (u UserRepository) GetUserById(db *sql.DB, user models.User, userId int) (models.User, error) {
	row := db.QueryRow("SELECT id, email, activated, username FROM users WHERE id = $1", userId)
	err := row.Scan(&user.ID, &user.Email, &user.Activated, &user.Username)

	if err != nil {

		return user, err
	}

	return user, nil
}

// CheckKey is for activating accounts, only called once per user
func (u UserRepository) CheckKey(db *sql.DB, id int, uuid string) error {
	var user models.User
	row := db.QueryRow("SELECT id, activation_key, activated FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.SignupKey, &user.Activated)

	if err != nil {
		log.Println(err)
		return errors.New("Something went wrong, please try again later.")
	}

	if user.Activated == true {
		log.Println("Already activated")
		return errors.New("This account is already activated")
	}

	if uuid == user.SignupKey && id == user.ID {
		db.Query("UPDATE users SET activated = TRUE WHERE id = $1", id)
		return nil
	} else {
		log.Println("There was an error verifying the email.")
		log.Println("Here are the values:")
		log.Println("user id that was provided:", id)
		log.Println("uuid that was passed in:", uuid)
		log.Println("user id from DB", user.ID)
		log.Println("Signup key from DB:", user.SignupKey)

		return errors.New("Something didn't match")
	}
}
