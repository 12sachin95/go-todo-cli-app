package services

import (
	"context"
	"go-todo-cli/db"
	"go-todo-cli/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddTodo adds a new todo to the MongoDB
func AddTodo(todo models.Todo) (*mongo.InsertOneResult, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	todo.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, todo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result, nil
}

// GetTodos retrieves all todos from MongoDB
func GetTodos() ([]models.Todo, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []models.Todo
	for cursor.Next(ctx) {
		var todo models.Todo
		err := cursor.Decode(&todo)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

// GetTodoByID retrieves a todo by its ID
func GetTodoByID(id string) (models.Todo, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	var todo models.Todo
	err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&todo)
	if err != nil {
		return todo, err
	}
	return todo, nil
}

// UpdateTodo updates an existing todo in MongoDB
func UpdateTodo(id string, updatedTodo models.Todo) (*mongo.UpdateResult, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"title": updatedTodo.Title}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteTodo deletes a todo by its ID from MongoDB
func DeleteTodo(id string) (*mongo.DeleteResult, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}
	return result, nil
}
