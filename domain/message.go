package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a message sent by a user.
type Message struct {
	MessageID primitive.ObjectID `json:"message_id" bson:"message_id,omitempty"`
	SenderID  primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	Content   string             `json:"content" bson:"content"`
	Time      time.Time          `json:"time" bson:"time"`
	Edited    bool               `json:"edited" bson:"edited"`
}

type MessageRepository interface {
	SendMessage(ctx context.Context, chatID primitive.ObjectID, message *Message) error
	GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]Message, error)
	GetMessage(ctx context.Context, chatID, messageID primitive.ObjectID) (Message, error) 
	DeleteMessage(ctx context.Context, chatID, messageID primitive.ObjectID) error
	UpdateMessage(ctx context.Context, chatID, messageID primitive.ObjectID, newContent string) error
}

type MessageUsecase interface {
	SendMessage(ctx context.Context, chatID primitive.ObjectID, message *Message) error
	GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]Message, error)
	GetMessage(ctx context.Context, chatID, messageID primitive.ObjectID) (Message, error) 
	DeleteMessage(ctx context.Context, chatID, messageID primitive.ObjectID) error
	UpdateMessage(ctx context.Context, chatID, messageID primitive.ObjectID, newContent string) error
}
