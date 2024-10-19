package cmd

import (
	"fmt"
	"go-todo-cli/api"
	"go-todo-cli/db"
	"log"

	"github.com/spf13/cobra"
)

// serveCmd starts the Gin server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		api.StartServer() // Start the Gin server
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	// Register the getToken command
	RootCmd.AddCommand(getTokenCmd)

	// Define the --user_id flag
	getTokenCmd.Flags().String("user_id", "", "User ID to get the token for")
	getTokenCmd.MarkFlagRequired("user_id")
}

// Define the 'getToken' command
var getTokenCmd = &cobra.Command{
	Use:   "getToken",
	Short: "Get a token for a specific user",
	Run: func(cmd *cobra.Command, args []string) {
		userID, err := cmd.Flags().GetString("user_id")
		if err != nil {
			log.Fatalf("Error reading user_id flag: %v", err)
		}

		token, err := db.GetTokenByUserID(userID)
		if err != nil {
			log.Fatalf("Error retrieving token: %v", err)
		}

		fmt.Printf("Token for user %s: %s\n", userID, token)
	},
}
