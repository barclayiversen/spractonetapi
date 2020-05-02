package models

type Data struct {
	Errs     map[string]string
	Email    string `validate:"required,email`
	Password string
	Username string
	Verified bool
}
