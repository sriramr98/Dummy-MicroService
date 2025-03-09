package services

import (
	"database/sql"
	"errors"
)

type AuthService struct {
	db *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Blocked  bool
}

func (u User) Validate() map[string]string {
	validationErrors := make(map[string]string)

	if u.Email == "" {
		validationErrors["email"] = "Email cannot be empty"
	}

	if u.Password == "" {
		validationErrors["password"] = "Password cannot be empty"
	}

	return validationErrors
}

var ErrUserExists = errors.New("user already exists")

func NewAuthService(db *sql.DB) AuthService {
	return AuthService{db: db}
}

func (as AuthService) CreateNewUser(email string, password string) (User, error) {
	_, err := as.GetUserInfo(email)
	if err == nil {
		// User already exists
		return User{}, ErrUserExists
	}

	user := User{
		Email:    email,
		Password: password,
	}

	err = as.db.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, blocked", email, password).Scan(&user.ID, &user.Blocked)
	return user, err
}

func (as AuthService) GetUserInfo(email string) (User, error) {
	var user User
	err := as.db.QueryRow("SELECT id, email, password, blocked FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password, &user.Blocked)
	return user, err
}
