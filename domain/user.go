package domain

import (
	"context"
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

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (primitive.ObjectID, error)
	GetUserByID(ctx context.Context, userID primitive.ObjectID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, userID primitive.ObjectID, user *User) error
	DeleteUser(ctx context.Context, userID primitive.ObjectID) error
}