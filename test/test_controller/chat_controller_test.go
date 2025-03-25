package test

import (
	"Real-Time-Chat-Application/controller"
	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/test/test_usecase/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateChat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Chat Created", func(t *testing.T) {
		mockChatUsecase := new(mocks.MockChatUsecase)
		chatController := controller.NewChatController(mockChatUsecase)

		r := gin.Default()
		r.POST("/chats", chatController.CreateChat)

		chat := domain.Chat{
			Participants: []primitive.ObjectID{primitive.NewObjectID()},
		}

		mockChatUsecase.On("CreateChat", mock.Anything, mock.AnythingOfType("*domain.Chat")).Return(primitive.NewObjectID(), nil)

		jsonValue, _ := json.Marshal(chat)
		req, _ := http.NewRequest("POST", "/chats", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockChatUsecase.AssertExpectations(t)
	})

	t.Run("Failure - Internal Server Error", func(t *testing.T) {
		mockChatUsecase := new(mocks.MockChatUsecase)
		chatController := controller.NewChatController(mockChatUsecase)

		r := gin.Default()
		r.POST("/chats", chatController.CreateChat)

		chat := domain.Chat{
			Participants: []primitive.ObjectID{primitive.NewObjectID()},
		}

		mockChatUsecase.On("CreateChat", mock.Anything, mock.AnythingOfType("*domain.Chat")).Return(primitive.NilObjectID, errors.New("failed to create chat"))

		jsonValue, _ := json.Marshal(chat)
		req, _ := http.NewRequest("POST", "/chats", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockChatUsecase.AssertExpectations(t)
	})

	t.Run("Failure - Invalid Request Body", func(t *testing.T) {
		mockChatUsecase := new(mocks.MockChatUsecase)
		chatController := controller.NewChatController(mockChatUsecase)

		r := gin.Default()
		r.POST("/chats", chatController.CreateChat)

		invalidJSON := `{"Participants": "invalid_data"}` // This should be an array of ObjectIDs

		req, _ := http.NewRequest("POST", "/chats", bytes.NewBuffer([]byte(invalidJSON)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockChatUsecase.AssertExpectations(t)
	})
}

func TestGetChat(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	chatController := controller.NewChatController(mockChatUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/chats/:chat_id", chatController.GetChat)

	chatID := primitive.NewObjectID()
	expectedChat := &domain.Chat{
		ChatID:       chatID,
		Participants: []primitive.ObjectID{primitive.NewObjectID()},
	}

	mockChatUsecase.On("GetChat", mock.Anything, chatID).Return(expectedChat, nil)

	req, _ := http.NewRequest("GET", "/chats/"+chatID.Hex(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Chat
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedChat.ChatID, response.ChatID)
}

func TestGetChatsByUserID(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	chatController := controller.NewChatController(mockChatUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/users/:user_id/chats", chatController.GetChatsByUserID)

	userID := primitive.NewObjectID()
	expectedChats := []domain.Chat{
		{
			ChatID:       primitive.NewObjectID(),
			Participants: []primitive.ObjectID{userID},
		},
	}

	mockChatUsecase.On("GetChatsByUserID", mock.Anything, userID).Return(expectedChats, nil)

	req, _ := http.NewRequest("GET", "/users/"+userID.Hex()+"/chats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.Chat
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedChats), len(response))
}

func TestUpdateChat(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	chatController := controller.NewChatController(mockChatUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/chats/:chat_id", chatController.UpdateChat)

	chatID := primitive.NewObjectID()
	chat := domain.Chat{
		Participants: []primitive.ObjectID{primitive.NewObjectID()},
	}

	mockChatUsecase.On("UpdateChat", mock.Anything, chatID, mock.AnythingOfType("*domain.Chat")).Return(nil)

	jsonValue, _ := json.Marshal(chat)
	req, _ := http.NewRequest("PUT", "/chats/"+chatID.Hex(), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteChat(t *testing.T) {
	mockChatUsecase := new(mocks.MockChatUsecase)
	chatController := controller.NewChatController(mockChatUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/chats/:chat_id", chatController.DeleteChat)

	chatID := primitive.NewObjectID()

	mockChatUsecase.On("DeleteChat", mock.Anything, chatID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/chats/"+chatID.Hex(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
