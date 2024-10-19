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

// ExtractUserIDFromJWT extracts user_id from the JWT claims
func ExtractUserIDFromJWT(c *gin.Context) {
	var jwtSecret = []byte(os.Getenv("SECRET_KEY")) // Secret key used to sign the JWT

	// Extract the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}

	// The token is usually in the format "Bearer <token>", so we split it
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	tokenString := parts[1]

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token's signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to extract claims"})
		c.Abort()
		return
	}

	// Extract the user_id from the claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		c.Abort()
		return
	}

	// Pass the userID to the context so it can be accessed in other handlers
	c.Set("userID", userID)
	c.Next() // Pass control to the next handler
}
