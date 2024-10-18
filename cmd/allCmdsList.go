package cmd

import (
	"encoding/json"
	"fmt"
	"go-todo-cli/utils"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// Struct to hold the response from the login API
type LoginResponse struct {
	Token string `json:"token"`
}

var registerCmd = &cobra.Command{
	Use:   "register [username] [password]",
	Short: "Register a new user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		username := args[0]
		password := args[1]

		resp, err := client.R().
			SetBody(map[string]string{"username": username, "password": password}).
			Post("http://localhost:8080/register")

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
		client := resty.New()
		username := args[0]
		password := args[1]

		resp, err := client.R().
			SetBody(map[string]string{"username": username, "password": password}).
			Post("http://localhost:8080/login")

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if resp.StatusCode() == 200 {
			// Unmarshal the JSON response to get the token
			var loginResponse LoginResponse
			err = json.Unmarshal(resp.Body(), &loginResponse)
			if err != nil {
				fmt.Println("Error parsing response:", err)
				return
			}

			// Save the token to the file
			token := loginResponse.Token

			err := utils.SaveTokenToFile(token) // Save the token to the file
			if err != nil {
				fmt.Println("Error saving token:", err)
				return
			}
			fmt.Println("Logged in successfully. Token saved.")
			fmt.Printf("Logged in! Token: %s\n", string(token))
		} else {
			fmt.Println("Login failed:", resp.String())
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and clear the JWT token",
	Run: func(cmd *cobra.Command, args []string) {
		// Delete the token file
		err := utils.DeleteTokenFile()
		if err != nil {
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
		// Load the token from the file
		token, err := utils.LoadTokenFromFile()
		if err != nil {
			fmt.Println("You need to log in first.")
			return
		}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+token). // Set the token for authorization
			SetBody(args[0]).
			Post("http://localhost:8080/todos")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("User created:", resp.String())
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "getOne [id]",
	Short: "Get a user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load the token from the file
		token, err := utils.LoadTokenFromFile()
		if err != nil {
			fmt.Println("You need to log in first.")
			return
		}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+token).
			Get(fmt.Sprintf("http://localhost:8080/todos/%s", args[0]))
		fmt.Println(resp, err)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("User details:", resp.String())
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [id] [name]",
	Short: "Update a user by ID",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Load the token from the file
		token, err := utils.LoadTokenFromFile()
		if err != nil {
			fmt.Println("You need to log in first.")
			return
		}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+token).
			SetBody(args[1]).
			Put(fmt.Sprintf("http://localhost:8080/todos/%s", args[0]))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("User updated:", resp.String())
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load the token from the file
		token, err := utils.LoadTokenFromFile()
		if err != nil {
			fmt.Println("You need to log in first.")
			return
		}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+token).
			Delete(fmt.Sprintf("http://localhost:8080/todos/%s", args[0]))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("User deleted:", resp.String())
		}
	},
}

// getCmd represents the get command
var getCmdAll = &cobra.Command{
	Use:   "get",
	Short: "Fetch all todos",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the token from the file
		token, err := utils.LoadTokenFromFile()
		if err != nil {
			fmt.Println("You need to log in first.")
			return
		}
		// Create a Resty client
		client := resty.New()

		// Send GET request to the API server
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+token). // Set the token for authorization
			Get("http://localhost:8080/todos")

		if err != nil {
			fmt.Println("Error fetching todos:", err)
			return
		}

		// Print the response
		fmt.Println(string(resp.Body()))
	},
}

func init() {
	RootCmd.AddCommand(registerCmd) // Add register command
	RootCmd.AddCommand(loginCmd)    // Add login command
	RootCmd.AddCommand(logoutCmd)

	RootCmd.AddCommand(createCmd)
	RootCmd.AddCommand(getCmd)
	RootCmd.AddCommand(updateCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(getCmdAll)
}
