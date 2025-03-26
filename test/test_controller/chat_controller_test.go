package test

import (
	"Real-Time-Chat-Application/controller"
	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/test/test_usecase/mocks"
	websocket "Real-Time-Chat-Application/websocket"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateChat(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	mockHub := &websocket.Hub{} // Create mock hub
	chatID := primitive.NewObjectID()
	senderID := primitive.NewObjectID()
	receiverID := primitive.NewObjectID()

	mockChatUsecase.On("CreateChat", mock.Anything, senderID, receiverID).Return(chatID, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request body
	reqBody := map[string]string{
		"sender_id":   senderID.Hex(),
		"receiver_id": receiverID.Hex(),
	}
	jsonBody, _ := json.Marshal(reqBody)
	c.Request = httptest.NewRequest("POST", "/chats", bytes.NewBuffer(jsonBody))

	// Create chat controller with mock usecase and hub
	chatController := controller.NewChatController(mockChatUsecase, mockHub)

	// Test the handler
	chatController.CreateChat(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	mockChatUsecase.AssertExpectations(t)
}

func TestGetChat(t *testing.T) {
    mockChatUsecase := new(mocks.MockChatUsecase)
    mockHub := &websocket.Hub{} // Create mock hub
    chatID := primitive.NewObjectID()

    // Creating the expected chat
    expectedChat := &domain.Chat{
        ChatID:       chatID,
        Participants: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
        Messages:     []domain.Message{},
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    // Mock the GetChat method to return the expected chat
    mockChatUsecase.On("GetChat", mock.Anything, chatID).Return(expectedChat, nil)

    // Creating a response recorder and test context
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    // Create a test request
    req := httptest.NewRequest("GET", "/chats/"+chatID.Hex(), nil)
    c.Request = req

    // Set the params in the request context (chat_id)
    c.Params = []gin.Param{{Key: "chat_id", Value: chatID.Hex()}}

    // Create chat controller with mock usecase and hub
    chatController := controller.NewChatController(mockChatUsecase, mockHub)

    // Run the handler
    chatController.GetChat(c)

    // Assert that the status code is 200 OK
    assert.Equal(t, http.StatusOK, w.Code)

    // Assert that expectations for mockChatUsecase were met
    mockChatUsecase.AssertExpectations(t)
}


func TestGetChatsByUserID(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	mockHub := &websocket.Hub{} // Create mock hub
	userID := primitive.NewObjectID()
	expectedChats := []domain.Chat{
		{
			ChatID:       primitive.NewObjectID(),
			Participants: []primitive.ObjectID{userID, primitive.NewObjectID()},
			Messages:     []domain.Message{},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	mockChatUsecase.On("GetChatsByUserID", mock.Anything, userID).Return(expectedChats, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Create request with user_id param
	req := httptest.NewRequest("GET", "/chats/user/"+userID.Hex(), nil)
	c.Request = req
	c.Params = []gin.Param{{Key: "user_id", Value: userID.Hex()}}

	// Create chat controller with mock usecase and hub
	chatController := controller.NewChatController(mockChatUsecase, mockHub)

	// Test the handler
	chatController.GetUserChats(c)
	assert.Equal(t, http.StatusOK, w.Code)

	mockChatUsecase.AssertExpectations(t)
}

func TestUpdateChat(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	mockHub := &websocket.Hub{} // Create mock hub
	chatID := primitive.NewObjectID()
	updatedChat := &domain.Chat{
		ChatID:       chatID,
		Participants: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
		Messages:     []domain.Message{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create a copy of updatedChat for the mock expectation
	expectedChat := &domain.Chat{
		ChatID:       updatedChat.ChatID,
		Participants: updatedChat.Participants,
		Messages:     updatedChat.Messages,
		CreatedAt:    updatedChat.CreatedAt,
		UpdatedAt:    updatedChat.UpdatedAt,
	}

	mockChatUsecase.On("UpdateChat", mock.Anything, chatID, mock.MatchedBy(func(chat *domain.Chat) bool {
		return chat.ChatID == expectedChat.ChatID &&
			len(chat.Participants) == len(expectedChat.Participants) &&
			len(chat.Messages) == len(expectedChat.Messages)
	})).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "chat_id", Value: chatID.Hex()}}

	jsonBody, _ := json.Marshal(updatedChat)
	c.Request = httptest.NewRequest("PUT", "/chats/"+chatID.Hex(), bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Test the handler
	chatController := controller.NewChatController(mockChatUsecase, mockHub)
	chatController.UpdateChat(c)
	
	assert.Equal(t, http.StatusOK, w.Code)

	mockChatUsecase.AssertExpectations(t)
}

func TestDeleteChat(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	mockHub := &websocket.Hub{} // Create mock hub
	chatID := primitive.NewObjectID()

	mockChatUsecase.On("DeleteChat", mock.Anything, chatID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "chat_id", Value: chatID.Hex()}}

	// Create request
	req := httptest.NewRequest("DELETE", "/chats/"+chatID.Hex(), nil)
	c.Request = req

	// Test the handler
	chatController := controller.NewChatController(mockChatUsecase, mockHub)
	chatController.DeleteChat(c)

	assert.Equal(t, http.StatusOK, w.Code)

	mockChatUsecase.AssertExpectations(t)
}

func TestGetChatByParticipants(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	mockHub := &websocket.Hub{} // Create mock hub
	senderID := primitive.NewObjectID()
	receiverID := primitive.NewObjectID()
	expectedChat := &domain.Chat{
		ChatID:       primitive.NewObjectID(),
		Participants: []primitive.ObjectID{senderID, receiverID},
		Messages:     []domain.Message{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockChatUsecase.On("GetChatByParticipants", mock.Anything, senderID, receiverID).Return(expectedChat, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "sender_id", Value: senderID.Hex()},
		{Key: "receiver_id", Value: receiverID.Hex()},
	}

	// Create request
	req := httptest.NewRequest("GET", fmt.Sprintf("/chats/%s/%s", 
		senderID.Hex(), receiverID.Hex()), nil)
	c.Request = req

	// Test the handler
	chatController := controller.NewChatController(mockChatUsecase, mockHub)
	chatController.GetChatByParticipants(c)
	assert.Equal(t, http.StatusOK, w.Code)

	mockChatUsecase.AssertExpectations(t)
}
