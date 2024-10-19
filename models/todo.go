package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo represents a task
type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string             `bson:"title" json:"title" validate:"required,min=1,max=100"` // Required, min length 1, max length 100
	Completed bool               `bson:"completed" json:"completed"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"` // Required User ID
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// TodoUpdate struct is used to update todo items
type TodoUpdate struct {
	Title     string    `json:"title,omitempty"`     // String, optional
	Completed *bool     `json:"completed,omitempty"` // Pointer to bool, optional
	UpdatedAt time.Time `json:"updated_at"`
}
