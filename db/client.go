package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient is the MongoDB client
var MongoClient *mongo.Client

// ConnectMongoDB connects to MongoDB
func ConnectMongoDB(uri string) {
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	MongoClient = client
}

// GetCollection returns a MongoDB collection
func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return MongoClient.Database(databaseName).Collection(collectionName)
}

// TokenData holds the token data structure from MongoDB
type TokenData struct {
	UserID string `bson:"user_id"`
	Token  string `bson:"token"`
	Exp    int64  `bson:"exp"`
}

// GetMongoClient connects to MongoDB and returns the client
func GetMongoClient() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI") // e.g., mongodb://localhost:27017
	if uri == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable is not set")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetTokenByUserID retrieves the token for a given user ID from MongoDB
func GetTokenByUserID(userID string) (string, error) {
	client, err := GetMongoClient()
	if err != nil {
		return "", fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	defer client.Disconnect(context.TODO())

	collection := client.Database(os.Getenv("DATABASE_NAME")).Collection("tokens") // Replace "mydb" and "tokens" with your DB and collection names

	var result TokenData
	filter := bson.M{"user_id": userID}

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error finding token: %v", err)
	}

	// Check if the token is expired
	if time.Now().Unix() > result.Exp {
		return "", fmt.Errorf("token has expired")
	}
	return result.Token, nil
}
