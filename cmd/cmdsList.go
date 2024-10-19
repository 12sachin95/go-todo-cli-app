package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Struct to hold the response from the login API
// type LoginResponse struct {
// 	Token string `json:"token"`
// }

// func getTokenFromMongoDB(ctx context.Context, collection *mongo.Collection) (string, error) {
// 	// Construct a query to find the token document
// 	filter := bson.M{}

// 	// Find the first token document
// 	result := collection.FindOne(ctx, filter)
// 	if result.Err() != nil {
// 		return "", fmt.Errorf("error finding token: %w", result.Err())
// 	}

// 	// Check if the result is nil
// 	if result.Err() == mongo.ErrNoDocuments {
// 		return "", fmt.Errorf("no token found in MongoDB")
// 	}

// 	// Decode the token document
// 	var tokenDoc struct {
// 		Token string `bson:"token"`
// 	}
// 	err := result.Decode(&tokenDoc)
// 	if err != nil {
// 		return "", fmt.Errorf("error decoding token: %w", err)
// 	}

// 	return tokenDoc.Token, nil
// }

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
		client := resty.New()
		username := args[0]
		password := args[1]

		resp, err := client.R().
			SetBody(map[string]string{"username": username, "password": password}).
			Post(TODO_SERVER_PATH + "/login")

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if resp.StatusCode() == 200 {
			// Unmarshal the JSON response to get the token
			// var loginResponse LoginResponse
			// err = json.Unmarshal(resp.Body(), &loginResponse)
			// if err != nil {
			// 	fmt.Println("Error parsing response:", err)
			// 	return
			// }

			// // Save the token to the file
			// token := loginResponse.Token

			// err := utils.SaveTokenToFile(token) // Save the token to the file
			// if err != nil {
			// 	fmt.Println("Error saving token:", err)
			// 	return
			// }
			// fmt.Println("Logged in successfully. Token saved.")
			fmt.Printf("Logged in! Token: %s\n", resp.String())
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
		// err := utils.DeleteTokenFile()
		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatal("Error loading .env file")
		// }

		// // MongoDB URI and options
		// mongoURI := os.Getenv("MONGODB_URI")
		clientOptions := options.Client().ApplyURI(MONGODB_URI)

		// Create a MongoDB client
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Disconnect client at the end
		defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
			}
		}()

		// Connect to the database and collection
		collection := client.Database("go-todo-db").Collection("tokens")

		// Query the token
		var result struct {
			Token string `bson:"token"`
		}

		// Find one document
		err = collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new Resty client

		client2 := resty.New()

		resp, err := client2.R().
			SetHeader("Authorization", "Bearer "+result.Token).
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
		// // Load the token from the file
		// // token, err := utils.LoadTokenFromFile()
		// // if err != nil {
		// // 	fmt.Println("You need to log in first.")
		// // 	return
		// // }
		// db.ConnectMongoDB()
		// ctx := context.Background()
		// collection := db.GetCollection("go-todo-db", "user")
		// token, err := getTokenFromMongoDB(ctx, collection)
		// if err != nil {
		// 	fmt.Println("User created:", err)
		// 	return
		// }
		// Load the .env file
		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatal("Error loading .env file")
		// }

		// // MongoDB URI and options
		// mongoURI := os.Getenv("MONGODB_URI")
		clientOptions := options.Client().ApplyURI(MONGODB_URI)

		// Create a MongoDB client
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Disconnect client at the end
		defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
			}
		}()

		// Connect to the database and collection
		collection := client.Database("go-todo-db").Collection("tokens")

		// Query the token
		var result struct {
			Token string `bson:"token"`
		}

		// Find one document
		err = collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new Resty client

		client2 := resty.New()
		resp, err := client2.R().
			SetHeader("Authorization", "Bearer "+result.Token). // Set the token for authorization
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
		// Load the token from the file
		// token, err := utils.LoadTokenFromFile()
		// if err != nil {
		// 	fmt.Println("You need to log in first.")
		// 	return
		// }
		// client := resty.New()

		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatal("Error loading .env file")
		// }

		// // MongoDB URI and options
		// mongoURI := os.Getenv("MONGODB_URI")
		clientOptions := options.Client().ApplyURI(MONGODB_URI)

		// Create a MongoDB client
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Disconnect client at the end
		defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
			}
		}()

		// Connect to the database and collection
		collection := client.Database("go-todo-db").Collection("tokens")

		// Query the token
		var result struct {
			Token string `bson:"token"`
		}

		// Find one document
		err = collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new Resty client

		client2 := resty.New()
		resp, err := client2.R().
			SetHeader("Authorization", "Bearer "+result.Token).
			Get(fmt.Sprintf(TODO_SERVER_PATH+"/todos/%s", args[0]))
		fmt.Println(resp, err)
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
		// Load the token from the file
		// token, err := utils.LoadTokenFromFile()
		// if err != nil {
		// 	fmt.Println("You need to log in first.")
		// 	return
		// }
		// client := resty.New()

		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatal("Error loading .env file")
		// }

		// // MongoDB URI and options
		// mongoURI := os.Getenv("MONGODB_URI")
		clientOptions := options.Client().ApplyURI(MONGODB_URI)

		// Create a MongoDB client
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Disconnect client at the end
		defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
			}
		}()

		// Connect to the database and collection
		collection := client.Database("go-todo-db").Collection("tokens")

		// Query the token
		var result struct {
			Token string `bson:"token"`
		}

		// Find one document
		err = collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new Resty client

		client2 := resty.New()
		resp, err := client2.R().
			SetHeader("Authorization", "Bearer "+result.Token).
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
		// Load the token from the file
		// token, err := utils.LoadTokenFromFile()
		// if err != nil {
		// 	fmt.Println("You need to log in first.")
		// 	return
		// }
		// client := resty.New()

		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatal("Error loading .env file")
		// }

		// // MongoDB URI and options
		// mongoURI := os.Getenv("MONGODB_URI")
		clientOptions := options.Client().ApplyURI(MONGODB_URI)

		// Create a MongoDB client
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Disconnect client at the end
		defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
			}
		}()

		// Connect to the database and collection
		collection := client.Database("go-todo-db").Collection("tokens")

		// Query the token
		var result struct {
			Token string `bson:"token"`
		}

		// Find one document
		err = collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new Resty client

		client2 := resty.New()
		resp, err := client2.R().
			SetHeader("Authorization", "Bearer "+result.Token).
			Delete(fmt.Sprintf(TODO_SERVER_PATH+"/todos/%s", args[0]))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("TODO deleted:", resp.String())
		}
	},
}

// getCmd represents the get command
var getCmdAll = &cobra.Command{
	Use:   "get",
	Short: "Fetch all todos",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the token from the file
		// token, err := utils.LoadTokenFromFile()
		// if err != nil {
		// 	fmt.Println("You need to log in first.")
		// 	return
		// }
		// // Create a Resty client
		// client := resty.New()

		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatal("Error loading .env file")
		// }

		// // MongoDB URI and options
		// mongoURI := os.Getenv("MONGODB_URI")
		clientOptions := options.Client().ApplyURI(MONGODB_URI)

		// Create a MongoDB client
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Disconnect client at the end
		defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
			}
		}()

		// Connect to the database and collection
		collection := client.Database("go-todo-db").Collection("tokens")

		// Query the token
		var result struct {
			Token string `bson:"token"`
		}

		// Find one document
		err = collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new Resty client

		client2 := resty.New()

		// Send GET request to the API server
		resp, err := client2.R().
			SetHeader("Authorization", "Bearer "+result.Token). // Set the token for authorization
			Get(TODO_SERVER_PATH + "/todos")

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
