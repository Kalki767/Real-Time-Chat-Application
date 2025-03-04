package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user of the application.
type User struct {
	UserID   primitive.ObjectID   `json:"-" bson:"_id"`
	Email    string               `json:"email" bson:"email"`
	Username string               `json:"username" bson:"username"`
	Password string               `json:"-" bson:"password"`
	Chats    []primitive.ObjectID `json:"chats" bson:"chats"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}