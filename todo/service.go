package todo

import (
	"errors"
	"fmt"
)

// Todo struct
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var NewTodo struct {
	Title string `json:"title" binding:"required"` // Require the title field
}

// In-memory todo storage
var (
	todos  = []Todo{}
	nextID = 1 // Start IDs from 1
)

// GetTodos returns all todos
func GetTodos() []Todo {
	return todos
}

// GetTodoByID returns a todo by its ID
func GetTodoByID(id int) (*Todo, error) {
	for _, todo := range todos {
		if todo.ID == id {
			fmt.Println(todo.ID)
			return &todo, nil
		}
	}

	return nil, errors.New("todo not found")
}

// AddTodo adds a new todo to the list
func AddTodo(newTodo Todo) {
	newTodo.ID = nextID // Set the ID to the next available ID
	todos = append(todos, newTodo)
	nextID++ // Increment the ID for the next todo
}

// UpdateTodoByID updates a todo by its ID
func UpdateTodoByID(id int, updatedTodo Todo) error {
	for i, todo := range todos {
		if todo.ID == id {
			todos[i] = updatedTodo
			return nil
		}
	}
	return errors.New("todo not found")
}

// DeleteTodo removes a todo by its ID
func DeleteTodo(id int) error {
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return errors.New("todo not found")
}
