package api

import (
	"context"
	"fmt"
	"go-todo-cli/db"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// AuthMiddleware verifies the JWT token for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Optionally check if the token is in the MongoDB (e.g., for logout)
		collection := db.GetCollection("go-todo-db", "tokens")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var storedToken bson.M
		err = collection.FindOne(ctx, bson.M{"token": tokenString}).Decode(&storedToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			c.Abort()
			return
		}

		// Allow the request to proceed
		c.Next()
	}
}
