package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Todo represents a task
type Todo struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title string             `bson:"title" json:"title"`
}
