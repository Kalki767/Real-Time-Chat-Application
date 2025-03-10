package repository

import (
	"Real-Time-Chat-Application/domain"
	"context"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChatRepository is a struct that defines the ChatRepository type
// send message sends a message between two people who have already been chatting. It accepts the chatId and the message to send
func(chatrepo *ChatRepository) SendMessage(ctx context.Context, chatID primitive.ObjectID, message *domain.Message) error {

	collection := chatrepo.collection

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

// Get messages will retrieve all the available messages based on their time stamp from the specific chat required
func(chatrepo *ChatRepository) GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]domain.Message, error) {

	collection := chatrepo.collection

	//get the messages using the chatID
	var chat domain.Chat

	//if the required chat doesn't exist in the database return there is no message to retireve with
	err := collection.FindOne(ctx,bson.M{"_id":chatID}).Decode(&chat)
	if err != nil {
		return []domain.Message{}, fmt.Errorf("failed to fetch chat: %w", err)
	}

	sort.Slice(chat.Messages, func(i,j int) bool {
		return chat.Messages[i].Time.Before(chat.Messages[j].Time)
	})

	//otherwise return list of messages in sorted order
	return chat.Messages, nil

}

// Getmessage returns a specific message given the message id from a specific chat
func(chatrepo *ChatRepository) GetMessage(ctx context.Context, chatID, messageID primitive.ObjectID) (domain.Message, error) {

	collection := chatrepo.collection

	// declare the parameters to retrieve the specified message from the chat list
	filter := bson.M{"_id":chatID, "messages.messageID":messageID}
	projection := bson.M{"messages.$": 1}
	
	// declare a list that will be type of message
	var result []domain.Message

	// retrieve the specified message from the database
	err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return domain.Message{}, fmt.Errorf("failed to fetch chat: %w", err)
	}

	//if the message is found return the message
	return result[0], nil

}

func(chatrepo *ChatRepository) DeleteMessage(ctx context.Context, chatID, messageID primitive.ObjectID) error{

	collection := chatrepo.collection

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

	//if no entry was modified return an error message
	if result.ModifiedCount == 0 {
		return fmt.Errorf("message not found or already deleted")
	}

	return nil
	
}

func (chatrepo *ChatRepository) UpdateMessage(ctx context.Context, chatID, messageID primitive.ObjectID, newContent string) error {
	collection := chatrepo.collection

	// Define the update query
	update := bson.M{
		"$set": bson.M{
			"messages.$[elem].content": newContent, // Update the message content
			"messages.$[elem].edited":  true,       // Mark the message as 
			"updated_at": time.Now(),
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