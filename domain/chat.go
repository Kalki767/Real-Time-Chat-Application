package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


// Chat represents a chat between two users.
type Chat struct {
	ChatID     primitive.ObjectID `json:"chat_id" bson:"_id,omitempty"`
	Participants []primitive.ObjectID `json:"participants" bson:"participants"` // [SenderID, ReceiverID]
	Messages   []Message          `json:"messages" bson:"messages"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type ChatRepository interface {
	CreateChat(chat *Chat) (primitive.ObjectID, error)
	GetChat(chatID primitive.ObjectID) (*Chat, error)
	GetChatsByUserID(userID primitive.ObjectID) ([]Chat, error)
	UpdateChat(chatID primitive.ObjectID, chat *Chat) error
	DeleteChat(chatID primitive.ObjectID) error
}