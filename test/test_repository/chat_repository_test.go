// test/chat_repository_test.go
package test

import (
	"context"
	"testing"

	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/repository"
	"Real-Time-Chat-Application/test/mongo/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateChat(t *testing.T) {
	// Setup
	mockCollection := new(mocks.MockCollection)
	repo := repository.NewChatRepository(mockCollection)

	SenderID := primitive.NewObjectID()
	ReceiverID := primitive.NewObjectID()
	chatID := primitive.NewObjectID()

	// expectedChat := domain.Chat{
	// 	Participants: []primitive.ObjectID{SenderID, ReceiverID},
	// 	Messages:     []domain.Message{},
	// }

	// Mock InsertOne to return a successful result
	mockCollection.On("InsertOne", mock.Anything, mock.MatchedBy(func(chat domain.Chat) bool {
		// Verify the chat object has the expected participants
		return len(chat.Participants) == 2 &&
			chat.Participants[0] == SenderID &&
			chat.Participants[1] == ReceiverID
	})).Return(&mongo.InsertOneResult{InsertedID: chatID}, nil)

	// Execute
	insertedID, err := repo.CreateChat(context.TODO(), SenderID, ReceiverID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, chatID, insertedID)
	mockCollection.AssertExpectations(t)
}

// Add other tests...

func TestGetChat(t *testing.T) {
	// Setup
	mockCollection := new(mocks.MockCollection)
	mockSingleResult := new(mocks.MockSingleResult)
	repo := repository.NewChatRepository(mockCollection)

	chatID := primitive.NewObjectID()
	expectedChat := &domain.Chat{
		ChatID:       chatID,
		Participants: []primitive.ObjectID{primitive.NewObjectID()},
		Messages:     []domain.Message{},
	}

	// Mock Decode
	mockSingleResult.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*domain.Chat)
		*arg = *expectedChat
	})

	// Mock FindOne
	mockCollection.On("FindOne", mock.Anything, bson.M{"_id": chatID}).Return(mockSingleResult)

	// Execute
	chat, err := repo.GetChat(context.TODO(), chatID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedChat, chat)
	mockCollection.AssertExpectations(t)
	mockSingleResult.AssertExpectations(t)
}

func TestGetChatsByUserID(t *testing.T) {
	// Create mock collection and cursor
	mockCollection := new(mocks.MockCollection)
	mockCursor := new(mocks.MockCursor)

	// Create a sample user ID
	userID := primitive.NewObjectID()

	// Create a sample chat
	expectedChat := domain.Chat{
		ChatID:       primitive.NewObjectID(),
		Participants: []primitive.ObjectID{userID},
		Messages:     []domain.Message{},
	}

	// Set up expectations for Find to return the mock cursor
	mockCollection.On("Find", mock.Anything, bson.M{"participants": userID}).Return(mockCursor, nil)

	// Set up expectations for the mock cursor
	mockCursor.On("Next", mock.Anything).Return(true).Once()
	mockCursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*domain.Chat)
		*arg = expectedChat
	}).Return(nil)
	mockCursor.On("Next", mock.Anything).Return(false).Once()
	mockCursor.On("Close", mock.Anything).Return(nil)
	mockCursor.On("Err").Return(nil)

	// Create repository
	chatrepo := repository.NewChatRepository(mockCollection)

	// Call the function
	chats, err := chatrepo.GetChatsByUserID(context.TODO(), userID)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, chats, 1)
	assert.Equal(t, expectedChat, chats[0])

	// Verify expectations
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestUpdateChat(t *testing.T) {
	// Setup
	mockCollection := new(mocks.MockCollection)
	repo := repository.NewChatRepository(mockCollection)

	chatID := primitive.NewObjectID()
	chat := &domain.Chat{
		ChatID:       chatID,
		Participants: []primitive.ObjectID{primitive.NewObjectID()},
		Messages:     []domain.Message{},
	}

	// Create a mock response for UpdateOne (MongoDB returns *mongo.UpdateResult)
	updateResult := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}

	// Mock UpdateOne to return the updateResult and nil error
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": chatID}, bson.M{"$set": chat}).Return(updateResult, nil)

	// Execute
	err := repo.UpdateChat(context.TODO(), chatID, chat)

	// Verify
	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
}

func TestDeleteChat(t *testing.T) {
	//set up 
	mockCollection := new(mocks.MockCollection)
	repo := repository.NewChatRepository(mockCollection)

	chatID := primitive.NewObjectID()

	// Create a mock response for DeleteOne (MongoDB returns *mongo.DeleteResult)
	deleteResult := &mongo.DeleteResult{DeletedCount: 1}

	// Mock DeleteOne to return the deleteResult and nil error
	mockCollection.On("DeleteOne", mock.Anything, bson.M{"_id": chatID}).Return(deleteResult, nil)

	//Execute
	err := repo.DeleteChat(context.TODO(), chatID)

	//Verify
	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
}

func TestGetChatByParticipants(t *testing.T) {
	// Setup
	mockCollection := new(mocks.MockCollection)
	mockSingleResult := new(mocks.MockSingleResult)
	repo := repository.NewChatRepository(mockCollection)

	SenderID := primitive.NewObjectID()
	ReceiverID := primitive.NewObjectID()
	chatID := primitive.NewObjectID()

	// Create a mock chat
	expectedChat := &domain.Chat{
		ChatID:       chatID,
		Participants: []primitive.ObjectID{SenderID, ReceiverID},
		Messages:     []domain.Message{},
	}

	// Mock FindOne
	mockCollection.On("FindOne", mock.Anything, bson.M{
		"participants": bson.M{
			"$all":  []primitive.ObjectID{SenderID, ReceiverID},
			"$size": 2,
		},
	}).Return(mockSingleResult)

	// Mock Decode to return the expected chat
	mockSingleResult.On("Decode", mock.AnythingOfType("*domain.Chat")).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*domain.Chat)
		*chat = *expectedChat
	}).Return(nil)

	// Execute
	chat, err := repo.GetChatByParticipants(context.TODO(), SenderID, ReceiverID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedChat, chat)
	mockCollection.AssertExpectations(t)
	mockSingleResult.AssertExpectations(t)
}

