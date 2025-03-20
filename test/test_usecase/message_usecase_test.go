package test_usecase

import (
	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/test/test_repository/mocks"
	"Real-Time-Chat-Application/usecase"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendMessage(t *testing.T) {
	// Setup
	mockMessageRepo := new(mocks.MockMessageRepository)
	messageUsecase := usecase.NewMessageUsecase(mockMessageRepo, 1*time.Second)

	chatID := primitive.NewObjectID()
	message := &domain.Message{
		MessageID: primitive.NewObjectID(),
		Content:   "Hello, world!",
		Time:      time.Now(),
	}

	// Mock the repository layer
	mockMessageRepo.On("SendMessage", mock.Anything, chatID, message).Return(nil)

	// Call the usecase layer
	err := messageUsecase.SendMessage(context.Background(), chatID, message)

	// Assert
	assert.NoError(t, err)
	mockMessageRepo.AssertExpectations(t)
}

func TestGetMessages(t *testing.T) {
	// Setup
	mockMessageRepo := new(mocks.MockMessageRepository)
	messageUsecase := usecase.NewMessageUsecase(mockMessageRepo, 1*time.Second)

	chatID := primitive.NewObjectID()
	messages := []domain.Message{
		{
			MessageID: primitive.NewObjectID(),
			Content:   "Hello, world!",
			Time:      time.Now(),
		},
		{
			MessageID: primitive.NewObjectID(),
			Content:   "Hello, world!",
			Time:      time.Now(),
		},
	}

	// Mock the repository layer
	mockMessageRepo.On("GetMessages", mock.Anything, chatID).Return(messages, nil)

	// Call the usecase layer
	receivedMessages, err := messageUsecase.GetMessages(context.Background(), chatID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, messages, receivedMessages)
	mockMessageRepo.AssertExpectations(t)
}

func TestGetMessage(t *testing.T) {
	// Setup
	mockMessageRepo := new(mocks.MockMessageRepository)
	messageUsecase := usecase.NewMessageUsecase(mockMessageRepo, 1*time.Second)

	chatID := primitive.NewObjectID()
	messageID := primitive.NewObjectID()
	message := domain.Message{
		MessageID: messageID,
		Content:   "Hello, world!",
		Time:      time.Now(),
	}

	// Mock the repository layer
	mockMessageRepo.On("GetMessage", mock.Anything, chatID, messageID).Return(message, nil)

	// Call the usecase layer
	receivedMessage, err := messageUsecase.GetMessage(context.Background(), chatID, messageID)	

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, message, receivedMessage)
	mockMessageRepo.AssertExpectations(t)
}

func TestUpdateMessage(t *testing.T) {
	// Setup
	mockMessageRepo := new(mocks.MockMessageRepository)
	messageUsecase := usecase.NewMessageUsecase(mockMessageRepo, 1*time.Second)

	chatID := primitive.NewObjectID()
	messageID := primitive.NewObjectID()	
	newContent := "Updated message"

	// Mock the repository layer
	mockMessageRepo.On("UpdateMessage", mock.Anything, chatID, messageID, newContent).Return(nil)

	// Call the usecase layer
	err := messageUsecase.UpdateMessage(context.Background(), chatID, messageID, newContent)	

	// Assert
	assert.NoError(t, err)
	mockMessageRepo.AssertExpectations(t)
}


func TestDeleteMessage(t *testing.T) {
	// Setup
	mockMessageRepo := new(mocks.MockMessageRepository)
	messageUsecase := usecase.NewMessageUsecase(mockMessageRepo, 1*time.Second)	

	chatID := primitive.NewObjectID()
	messageID := primitive.NewObjectID()

	// Mock the repository layer
	mockMessageRepo.On("DeleteMessage", mock.Anything, chatID, messageID).Return(nil)	

	// Call the usecase layer
	err := messageUsecase.DeleteMessage(context.Background(), chatID, messageID)

	// Assert
	assert.NoError(t, err)
	mockMessageRepo.AssertExpectations(t)
}	

