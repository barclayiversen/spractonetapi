package userRepository

import (
	"database/sql"
	"fmt"
	"log"
	"spractonetapi/models"
)

type UserRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (u UserRepository) Signup(db *sql.DB, user models.User) (models.User, error) {
	stmt := "INSERT INTO USERS (email,password,username) VALUES ($1, $2, $3) RETURNING id;"
	err := db.QueryRow(stmt, user.Email, user.Password, user.Username).Scan(&user.ID)

	if err != nil {
		fmt.Println(err)
		return user, err
	}

	user.Password = ""
	return user, nil
}

func (u UserRepository) Login(db *sql.DB, user models.User) (models.User, error) {

	row := db.QueryRow("SELECT id, email, username, password FROM users WHERE email = $1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u UserRepository) GetUserById(db *sql.DB, user models.User, userId int) (models.User, error) {
	row := db.QueryRow("SELECT id, email, username FROM users WHERE id = $1", userId)
	err := row.Scan(&user.ID, &user.Email, &user.Username)

	if err != nil {

		return user, err
	}

	return user, nil
}
