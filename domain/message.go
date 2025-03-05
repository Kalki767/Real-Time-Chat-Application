package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a message sent by a user.
type Message struct {
	MessageID primitive.ObjectID `json:"message_id" bson:"message_id,omitempty"`
	SenderID  primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	Content   string             `json:"content" bson:"content"`
	Time      time.Time          `json:"time" bson:"time"`
}

type MessageRepository interface {
	SendMessage(chatID primitive.ObjectID, message *Message) error
	GetMessages(chatID primitive.ObjectID) ([]Message, error)
	GetMessage(chatID, messageID primitive.ObjectID) (*Message, error) 
	DeleteMessage(chatID, messageID primitive.ObjectID) error
}