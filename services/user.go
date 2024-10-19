package services

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
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Token string `bson:"token"`
	// Add other fields as needed (e.g., expiration time, user ID)
}

// func generateToken(claims jwt.MapClaims) (string, error) {
// 	// Set token expiration time
// 	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Adjust expiration as needed

// 	// Create a new JWT token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// Sign the token with a secret key
// 	return token.SignedString([]byte("your-secret-key")) // Replace with a secure secret
// }
// func saveTokenToMongoDB(ctx context.Context, collection *mongo.Collection, token string) error {
// 	// Create a token document
// 	tokenDoc := Token{
// 		Token: token,
// 		// Add other fields if needed
// 	}

// 	// Insert the token document into the collection
// 	_, err := collection.InsertOne(ctx, tokenDoc)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// RegisterUser adds a new user to MongoDB
func RegisterUser(username, password string) (*mongo.InsertOneResult, error) {
	collection := db.GetCollection("go-todo-db", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Hash the password before storing
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: string(passwordHash),
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AuthenticateUser authenticates a user and returns a JWT token
func AuthenticateUser(username, password string) (string, error) {
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

// LogoutUser invalidates the JWT token by deleting it from MongoDB
func LogoutUser(tokenString string) error {
	collection := db.GetCollection("go-todo-db", "tokens")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"token": tokenString})
	return err
}
