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

	result, err := collection.InsertOne(ctx, todo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result, nil
}

// GetTodos retrieves all todos from MongoDB
func GetTodos(userId primitive.ObjectID) ([]models.Todo, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	cursor, err := collection.Find(ctx, filter)
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
func GetTodoByID(id string, userId primitive.ObjectID) (models.Todo, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	var todo models.Todo
	err := collection.FindOne(ctx, bson.M{"_id": objectID, "user_id": userId}).Decode(&todo)
	if err != nil {
		return todo, err
	}
	return todo, nil
}

// UpdateTodo updates an existing todo in MongoDB
func UpdateTodo(id string, userId primitive.ObjectID, updatedTodo models.TodoUpdate) (*mongo.UpdateResult, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID, "user_id": userId}
	updateFields := bson.M{
		"updated_at": updatedTodo.UpdatedAt,
	}
	if updatedTodo.Title != "" {
		updateFields["title"] = updatedTodo.Title
	}
	if updatedTodo.Completed != nil {
		updateFields["completed"] = *updatedTodo.Completed // Dereference the pointer to get the actual bool value
	}

	update := bson.M{"$set": updateFields}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteTodo deletes a todo by its ID from MongoDB
func DeleteTodo(id string, userId primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := db.GetCollection("go-todo-db", "todos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(id)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID, "user_id": userId})
	if err != nil {
		return nil, err
	}
	return result, nil
}
