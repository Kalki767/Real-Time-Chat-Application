package repository

import (
	"Real-Time-Chat-Application/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func(chatrepo *ChatRepository) SendMessage(ctx context.Context, chatID primitive.ObjectID, message *domain.Message) error {

	collection := chatrepo.database.Collection(chatrepo.collection)

	// add a time stamp to the message
	message.MessageID = primitive.NewObjectID()
	message.Time = time.Now()

	// Push the message into the "messages" array in the chat document
	update := bson.M{
		"$push": bson.M{"messages": message},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	// add the message to the chat list inside the chat ID
	_, err := collection.UpdateOne(ctx, bson.M{"_id": chatID}, update)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil

}

func(chatrepo *ChatRepository) GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]domain.Message, error) {

	collection := chatrepo.database.Collection(chatrepo.collection)

	//get the messages using the chatID
	var chat domain.Chat

	err := collection.FindOne(ctx,bson.M{"_id":chatID}).Decode(&chat)
	if err != nil {
		return chat.Messages, fmt.Errorf("failed to fetch chat: %w", err)
	}

	return chat.Messages, nil

}

func(chatrepo *ChatRepository) GetMessage(ctx context.Context, chatID, messageID primitive.ObjectID) (domain.Message, error) {

	collection := chatrepo.database.Collection(chatrepo.collection)

	//get the messages using the chatID
	var chat domain.Chat

	var message domain.Message
	err := collection.FindOne(ctx,bson.M{"_id":chatID}).Decode(&chat)
	if err != nil {
		return message, fmt.Errorf("failed to fetch chat: %w", err)
	}

	

	for _, msg := range chat.Messages {
		if msg.MessageID == messageID {
			return msg, nil
		}
	}
	return message, fmt.Errorf("message not found")

}

func(chatrepo *ChatRepository) DeleteMessage(ctx context.Context, chatID, messageID primitive.ObjectID) error{

	collection := chatrepo.database.Collection(chatrepo.collection)

	// Pull the message from the "messages" array in the chat document
	update := bson.M{
		"$pull": bson.M{"messages": bson.M{"message_id": messageID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	// delete the message from the chat list inside the chat ID
	result, err := collection.UpdateOne(ctx,bson.M{"_id": chatID}, update)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("message not found or already deleted")
	}

	return nil
	
}

func (chatrepo *ChatRepository) UpdateMessage(ctx context.Context, chatID, messageID primitive.ObjectID, newContent string) error {
	collection := chatrepo.database.Collection(chatrepo.collection)

	// Define the update query
	update := bson.M{
		"$set": bson.M{
			"messages.$[elem].content": newContent, // Update the message content
			"messages.$[elem].edited":  true,       // Mark the message as edited
		},
	}

	// Define the filter to match the chat and the specific message inside the messages array
	filter := bson.M{"_id": chatID}
	arrayFilter := options.ArrayFilters{Filters: bson.A{bson.M{"elem.message_id": messageID}}}

	// Execute the update
	result, err := collection.UpdateOne(ctx, filter, update, &options.UpdateOptions{ArrayFilters: &arrayFilter})
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	// Check if any document was actually modified
	if result.ModifiedCount == 0 {
		return fmt.Errorf("message not found or already updated")
	}

	return nil
}