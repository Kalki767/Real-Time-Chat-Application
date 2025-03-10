package repository

import (
	"Real-Time-Chat-Application/domain"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type ChatRepository struct {
	collection *mongo.Collection
}

func NewChatRepository (database *mongo.Database, collection string) domain.ChatRepository {
	repo := &ChatRepository{collection: database.Collection(collection)}

	// Ensure an index on message_id for faster lookups
	coll := database.Collection(collection)
	indexModel := mongo.IndexModel{
		Keys: bson.M{"messages.message_id": 1},
	}
	_, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal("Failed to create index:", err)
	}

	return repo
}

func(chatrepo *ChatRepository) CreateChat(ctx context.Context,chat *domain.Chat) (primitive.ObjectID, error) {

	collection := chatrepo.collection

	result, err := collection.InsertOne(ctx, chat)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create chat: %w", err)
	}
	return result.InsertedID.(primitive.ObjectID), nil

}
func(chatrepo *ChatRepository) GetChat(ctx context.Context, chatID primitive.ObjectID) (*domain.Chat, error) {

	collection := chatrepo.collection
	var chat domain.Chat
	err := collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, fmt.Errorf("failed to fetch chat: %w", err)
	}
	return &chat, nil

}
func(chatrepo *ChatRepository) GetChatsByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Chat, error) {

	collection := chatrepo.collection
	cursor, err := collection.Find(ctx, bson.M{"participants": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chats: %w", err)
	}
	defer cursor.Close(ctx)

	var chats []domain.Chat
	for cursor.Next(ctx) {
		var chat domain.Chat
		if err := cursor.Decode(&chat); err != nil {
			return nil, fmt.Errorf("failed to decode chat: %w", err)
		}
		chats = append(chats, chat)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return chats, nil

}
func(chatrepo *ChatRepository) UpdateChat(ctx context.Context, chatID primitive.ObjectID, chat *domain.Chat) error {

	collection := chatrepo.collection
	_, err := collection.UpdateOne(ctx, bson.M{"_id": chatID}, bson.M{"$set": chat})
	if err != nil {
		return fmt.Errorf("failed to update chat: %w", err)
	}
	return nil
}
func(chatrepo *ChatRepository) DeleteChat(ctx context.Context, chatID primitive.ObjectID) error {

	collection := chatrepo.collection
	_, err := collection.DeleteOne(ctx, bson.M{"_id": chatID})
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}
	return nil

}