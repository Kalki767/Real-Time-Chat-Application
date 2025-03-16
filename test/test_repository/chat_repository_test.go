// test/chat_repository_test.go
package test

import (
	"context"
	"testing"

	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/repository"
	"Real-Time-Chat-Application/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateChat(t *testing.T) {
	// Setup
	mockCollection := new(mocks.MockCollection)
	repo := repository.NewChatRepository(mockCollection) // Pass the mock

	chat := &domain.Chat{
		ChatID:           primitive.NewObjectID(),
		Participants: []primitive.ObjectID{primitive.NewObjectID()},
		Messages:     []domain.Message{},
	}

	// Mock InsertOne to return a successful result
	mockCollection.On("InsertOne", mock.Anything, chat).Return(&mongo.InsertOneResult{InsertedID: chat.ChatID}, nil)

	// Execute
	insertedID, err := repo.CreateChat(context.TODO(), chat)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, chat.ChatID, insertedID)
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
	// Setup
	mockCollection := new(mocks.MockCollection)
	cursor := new(mocks.MockCursor)
	repo := repository.NewChatRepository(mockCollection)

	userID := primitive.NewObjectID()
	chatID := primitive.NewObjectID()
	expectedChats := []domain.Chat{
		{
			ChatID:       chatID,
			Participants: []primitive.ObjectID{userID},
			Messages:     []domain.Message{},
		},
	}

	// Mock Find
	mockCollection.On("Find", mock.Anything, bson.M{"participants": userID}).Return(cursor, nil)

	// Mock cursor
	cursor.On("Next", mock.Anything).Return(true).Once()
	cursor.On("Next", mock.Anything).Return(false).Once()
	cursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*domain.Chat)
		*arg = expectedChats[0]
	}).Return(nil).Once()
	cursor.On("Close", mock.Anything).Return(nil).Once()

	// Execute
	chats, err := repo.GetChatsByUserID(context.TODO(), userID)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedChats, chats)
	mockCollection.AssertExpectations(t)
	cursor.AssertExpectations(t)
}
