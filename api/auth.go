package api

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User struct
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User store
var users = make(map[string]string) // Simple map to store users

// Secret key for signing JWT
var secretKey = []byte("your_secret_key")

// GenerateJWT generates a JWT token
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Register a new user
func Register(username, password string) error {
	if _, exists := users[username]; exists {
		return errors.New("user already exists")
	}
	users[username] = password
	return nil
}

// Authenticate user credentials and return a JWT token
func Authenticate(username, password string) (string, error) {
	if storedPassword, exists := users[username]; exists && storedPassword == password {
		return GenerateJWT(username)
	}
	return "", errors.New("invalid credentials")
}
