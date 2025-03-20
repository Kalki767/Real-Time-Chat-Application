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

func TestCreateChat(t *testing.T) {
	mockChatRepository := new(mocks.MockChatRepository)
	chatUsecase := usecase.NewChatUsecase(mockChatRepository, 1*time.Second)

	chat := domain.Chat{
		ChatID:       primitive.NewObjectID(),
		Participants: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
		Messages:     []domain.Message{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockChatRepository.On("CreateChat", mock.Anything, &chat).Return(chat.ChatID, nil)

	chatID, err := chatUsecase.CreateChat(context.Background(), &chat)
	assert.NoError(t, err)
	assert.Equal(t, chat.ChatID, chatID)
	mockChatRepository.AssertExpectations(t)
}

func TestGetChat(t *testing.T) {
	mockChatRepository := new(mocks.MockChatRepository)
	chatUsecase := usecase.NewChatUsecase(mockChatRepository, 1*time.Second)

	chatID := primitive.NewObjectID()
	expectedchat := domain.Chat{
		ChatID:       chatID,
		Participants: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
		Messages:     []domain.Message{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockChatRepository.On("GetChat", mock.Anything, chatID).Return(&expectedchat, nil)

	chat, err := chatUsecase.GetChat(context.Background(), chatID)
	assert.NoError(t, err)
	assert.Equal(t, chat.ChatID, chatID)
	mockChatRepository.AssertExpectations(t)
}

func TestGetChatsByUserID(t *testing.T) {
	mockChatRepository := new(mocks.MockChatRepository)
	chatUsecase := usecase.NewChatUsecase(mockChatRepository, 1*time.Second)

	userID := primitive.NewObjectID()	
	expectedchats := []domain.Chat{
		{
			ChatID:       primitive.NewObjectID(),
			Participants: []primitive.ObjectID{userID, primitive.NewObjectID()},
			Messages:     []domain.Message{},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	mockChatRepository.On("GetChatsByUserID", mock.Anything, userID).Return(expectedchats, nil)

	chats, err := chatUsecase.GetChatsByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedchats, chats)
	mockChatRepository.AssertExpectations(t)
}

func TestUpdateChat(t *testing.T) {
	mockChatRepository := new(mocks.MockChatRepository)
	chatUsecase := usecase.NewChatUsecase(mockChatRepository, 1*time.Second)

	chatID := primitive.NewObjectID()
	updatedChat := domain.Chat{
		ChatID:       chatID,
		Participants: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
		Messages:     []domain.Message{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockChatRepository.On("UpdateChat", mock.Anything, chatID, &updatedChat).Return(nil)

	err := chatUsecase.UpdateChat(context.Background(), chatID, &updatedChat)
	assert.NoError(t, err)
	mockChatRepository.AssertExpectations(t)
}

func TestDeleteChat(t *testing.T) {
	mockChatRepository := new(mocks.MockChatRepository)
	chatUsecase := usecase.NewChatUsecase(mockChatRepository, 1*time.Second)

	chatID := primitive.NewObjectID()
	mockChatRepository.On("DeleteChat", mock.Anything, chatID).Return(nil)

	err := chatUsecase.DeleteChat(context.Background(), chatID)
	assert.NoError(t, err)
	mockChatRepository.AssertExpectations(t)
}





