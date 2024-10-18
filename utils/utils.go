package utils

import (
	"os"
)

// Helper function to save token to a file
func SaveTokenToFile(token string) error {
	return os.WriteFile("token.txt", []byte(token), 0644)
}

// Helper function to load token from a file
func LoadTokenFromFile() (string, error) {
	data, err := os.ReadFile("token.txt")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Helper function to delete the token file
func DeleteTokenFile() error {
	return os.Remove("token.txt")
}
