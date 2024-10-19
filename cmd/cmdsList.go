package cmd

import (
	"fmt"
	"go-todo-cli/db"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// Group command: `todoCmd`
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Commands related to todos",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Ensure user_id flag is set for todo commands
		if userID == "" {
			log.Fatalf("The --user_id flag is required for this command.")
		}
	},
}

// Group command: `userCmd`
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Commands related to users",
}

// Declare userID at the package level
var userID string

func init() {

	// add the groups of comnads
	RootCmd.AddCommand(userCmd)
	RootCmd.AddCommand(todoCmd)

	userCmd.AddCommand(registerCmd) // Add register command
	userCmd.AddCommand(loginCmd)    // Add login command

	// Define the --user_id flag as a persistent flag foronly logout cmd in user group
	logoutCmd.Flags().String("user_id", "", "User ID to get the token for")
	logoutCmd.MarkFlagRequired("user_id")
	userCmd.AddCommand(logoutCmd)

	// Register the getToken command
	userCmd.AddCommand(getTokenCmd)

	// Define the --user_id flag
	getTokenCmd.Flags().String("user_id", "", "User ID to get the token for")
	getTokenCmd.MarkFlagRequired("user_id")

	// Define the --user_id flag as a persistent flag for the `todoCmd` group
	todoCmd.PersistentFlags().StringVar(&userID, "user_id", "", "User ID to perform actions on todos")
	todoCmd.MarkPersistentFlagRequired("user_id")
	todoCmd.AddCommand(createCmd)
	todoCmd.AddCommand(getCmd)
	todoCmd.AddCommand(updateCmd)
	todoCmd.AddCommand(deleteCmd)
	todoCmd.AddCommand(getAllCmd)
}

// GetTokenForUser retrieves the token for a given user_id from the command flags
func GetTokenForUser(cmd *cobra.Command) (string, error) {
	// Get the user_id from the command flag
	userID, err := cmd.Flags().GetString("user_id")
	if err != nil {
		log.Fatalf("Error reading user_id flag: %v", err)
	}

	// Fetch the token from MongoDB using the user_id
	token, err := db.GetTokenByUserID(userID)
	if err != nil {
		log.Fatalf("Error retrieving token for user_id %s: %v", userID, err)
	}

	return token, nil
}

var registerCmd = &cobra.Command{
	Use:   "register [username] [password]",
	Short: "Register a new user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		restyClient := resty.New()
		username := args[0]
		password := args[1]

		resp, err := restyClient.R().
			SetBody(map[string]string{"username": username, "password": password}).
			Post(TODO_SERVER_PATH + "/register")

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if resp.StatusCode() == 201 {
			fmt.Println("User registered successfully.")
		} else {
			fmt.Println("Registration failed:", resp.String())
		}
	},
}

var loginCmd = &cobra.Command{
	Use:   "login [username] [password]",
	Short: "Login and get a JWT token",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		restyClient := resty.New()
		username := args[0]
		password := args[1]

		resp, err := restyClient.R().
			SetBody(map[string]string{"username": username, "password": password}).
			Post(TODO_SERVER_PATH + "/login")

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if resp.StatusCode() == 200 {
			fmt.Printf("Logged in! ")
		} else {
			fmt.Println("Login failed:", resp.String())
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and clear the JWT token",
	Run: func(cmd *cobra.Command, args []string) {
		// Use the reusable function to get the token
		token, err := GetTokenForUser(cmd)
		if err != nil {
			log.Fatalf("Failed to get token: %v", err)
		}

		// Create a new Resty Client
		restyClient := resty.New()

		resp, err := restyClient.R().
			SetHeader("Authorization", "Bearer "+token).
			Post(TODO_SERVER_PATH + "/logout")

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if err != nil || resp.StatusCode() != 200 {
			fmt.Println("Error logging out:", err)
			return
		}
		fmt.Println("Logged out successfully.")
	},
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetTokenForUser(cmd)
		if err != nil {
			log.Fatalf("Failed to get token: %v", err)
		}

		// Create a new Resty Client
		restyClient := resty.New()
		resp, err := restyClient.R().
			SetHeader("Authorization", "Bearer "+token). // Set the token for authorization
			SetBody(args[0]).
			Post(TODO_SERVER_PATH + "/todos")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("TODO created:", resp.String())
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "getOne [id]",
	Short: "Get a user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetTokenForUser(cmd)
		if err != nil {
			log.Fatalf("Failed to get token: %v", err)
		}

		// Create a new Resty Client
		restyClient := resty.New()
		resp, err := restyClient.R().
			SetHeader("Authorization", "Bearer "+token).
			Get(fmt.Sprintf(TODO_SERVER_PATH+"/todos/%s", args[0]))

		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("TODO details:", resp.String())
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [id] [name]",
	Short: "Update a user by ID",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetTokenForUser(cmd)
		if err != nil {
			log.Fatalf("Failed to get token: %v", err)
		}

		// Create a new Resty Client
		restyClient := resty.New()
		resp, err := restyClient.R().
			SetHeader("Authorization", "Bearer "+token).
			SetBody(args[1]).
			Put(fmt.Sprintf(TODO_SERVER_PATH+"/todos/%s", args[0]))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("TODO updated:", resp.String())
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetTokenForUser(cmd)
		if err != nil {
			log.Fatalf("Failed to get token: %v", err)
		}

		// Create a new Resty Client
		restyClient := resty.New()
		resp, err := restyClient.R().
			SetHeader("Authorization", "Bearer "+token).
			Delete(fmt.Sprintf(TODO_SERVER_PATH+"/todos/%s", args[0]))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("TODO deleted:", resp.String())
		}
	},
}

// getCmd represents the get command
var getAllCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch all todos",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetTokenForUser(cmd)
		if err != nil {
			log.Fatalf("Failed to get token: %v", err)
		}

		// Create a new Resty Client
		restyClient := resty.New()
		// Send GET request to the API server
		resp, err := restyClient.R().
			SetHeader("Authorization", "Bearer "+token). // Set the token for authorization
			Get(TODO_SERVER_PATH + "/todos")

		if err != nil {
			fmt.Println("Error fetching todos:", err)
			return
		}

		// Print the response
		fmt.Println(string(resp.Body()))
	},
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
