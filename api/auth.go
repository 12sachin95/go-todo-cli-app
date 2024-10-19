package api

// import (
// 	"errors"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// )

// // User struct
// type User struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// // User store
// var users = make(map[string]string) // Simple map to store users

// // Secret key for signing JWT
// var secretKey = []byte("your_secret_key")

// // GenerateJWT generates a JWT token
// func GenerateJWT(username string) (string, error) {
// 	claims := jwt.MapClaims{
// 		"username": username,
// 		"exp":      time.Now().Add(time.Hour * 72).Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(secretKey)
// }

// // Register a new user
// func Register(username, password string) error {
// 	if _, exists := users[username]; exists {
// 		return errors.New("user already exists")
// 	}
// 	users[username] = password
// 	return nil
// }

// // Authenticate user credentials and return a JWT token
// func Authenticate(username, password string) (string, error) {
// 	if storedPassword, exists := users[username]; exists && storedPassword == password {
// 		return GenerateJWT(username)
// 	}
// 	return "", errors.New("invalid credentials")
// }

// package services

import (
	"context"
	"errors"
	"go-todo-cli/db"
	"go-todo-cli/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser adds a new user to MongoDB
func Register(username, password string) error {
	collection := db.GetCollection("go-todo-db", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if username already exists
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return errors.New("username already exists")
	}

	// Hash the password before storing
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: string(passwordHash),
	}

	_, err = collection.InsertOne(ctx, user)
	return err
}

// var jwtSecret = []byte("your_secret_key")

// AuthenticateUser authenticates a user and returns a JWT token
func Authenticate(username, password string) (string, error) {
	collection := db.GetCollection("go-todo-db", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// Store the token in MongoDB for future validation (optional, used for logout)
	tokensCollection := db.GetCollection("go-todo-db", "tokens")
	tokensCollection.InsertOne(ctx, bson.M{"token": tokenString})

	return tokenString, nil
}
