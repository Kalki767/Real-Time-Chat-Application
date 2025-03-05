// package repository

// import (
// 	"Real-Time-Chat-Application/domain"
// 	"context"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type MessageRepository struct {
// 	database   *mongo.Database
// 	collection string
// }


// func NewMessageRepository(database *mongo.Database, collection string) *MessageRepository {
// 	return &MessageRepository{database: database, collection: collection}
// }

// func(msg *MessageRepository) SendMessage(ctx context.Context, chatID primitive.ObjectID, message *domain.Message) error {

// }

// func(msg *MessageRepository) GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]domain.Message, error) {

// }

// func(msg *MessageRepository) GetMessage(ctx context.Context, chatID, messageID primitive.ObjectID) (*domain.Message, error) {

// }

// func(msg *MessageRepository) DeleteMessage(ctx context.Context, chatID, messageID primitive.ObjectID) error{
	
// }