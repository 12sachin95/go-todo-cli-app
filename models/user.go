package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username" validate:"required,min=3,max=32"`
	Email    string             `bson:"email" json:"email" validate:"required,email"`
	Password string             `bson:"password" json:"password" validate:"required,min=3"`
}
