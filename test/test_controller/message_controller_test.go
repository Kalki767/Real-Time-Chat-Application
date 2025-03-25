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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendMessage(t *testing.T) {
	mockMessageUsecase := new(mocks.MockMessageUsecase)
	messageController := controller.NewMessageController(mockMessageUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/chats/:chat_id/messages", messageController.SendMessage)

	t.Run("success", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		message := &domain.Message{
			SenderID: primitive.NewObjectID(),
			Content:  "Test message",
			Edited:   false,
		}

		mockMessageUsecase.On("SendMessage", mock.Anything, chatID, message).Return(nil).Once()

		jsonMessage, _ := json.Marshal(message)
		req, _ := http.NewRequest(http.MethodPost, "/chats/"+chatID.Hex()+"/messages", bytes.NewBuffer(jsonMessage))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})

	t.Run("invalid chat ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/chats/invalid-id/messages", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		req, _ := http.NewRequest(http.MethodPost, "/chats/"+chatID.Hex()+"/messages", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		message := &domain.Message{
			SenderID: primitive.NewObjectID(),
			Content:  "Test message",
			Edited:   false,
		}

		mockMessageUsecase.On("SendMessage", mock.Anything, chatID, message).Return(errors.New("usecase error")).Once()

		jsonMessage, _ := json.Marshal(message)
		req, _ := http.NewRequest(http.MethodPost, "/chats/"+chatID.Hex()+"/messages", bytes.NewBuffer(jsonMessage))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})
}


func TestGetMessages(t *testing.T) {
	mockMessageUsecase := new(mocks.MockMessageUsecase)
	messageController := controller.NewMessageController(mockMessageUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/chats/:chat_id/messages", messageController.GetMessages)

	t.Run("success", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		expectedMessages := []domain.Message{
			{
				MessageID: primitive.NewObjectID(),
				SenderID:  primitive.NewObjectID(),
				Content:   "Test message 1",
				Time:      time.Now(),
			},
			{
				MessageID: primitive.NewObjectID(),
				SenderID:  primitive.NewObjectID(),
				Content:   "Test message 2",
				Time:      time.Now(),
			},
		}

		mockMessageUsecase.On("GetMessages", mock.Anything, chatID).Return(expectedMessages, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/chats/"+chatID.Hex()+"/messages", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})

	t.Run("invalid chat ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/chats/invalid-id/messages", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetMessage(t *testing.T) {
	mockMessageUsecase := new(mocks.MockMessageUsecase)
	messageController := controller.NewMessageController(mockMessageUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/chats/:chat_id/messages/:message_id", messageController.GetMessage)

	t.Run("success", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		messageID := primitive.NewObjectID()
		expectedMessage := domain.Message{
			MessageID: messageID,
			SenderID:  primitive.NewObjectID(),
			Content:   "Test message",
			Time:      time.Now(),
		}

		mockMessageUsecase.On("GetMessage", mock.Anything, chatID, messageID).Return(expectedMessage, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/chats/"+chatID.Hex()+"/messages/"+messageID.Hex(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})

	t.Run("invalid IDs", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/chats/invalid-id/messages/invalid-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteMessage(t *testing.T) {
	mockMessageUsecase := new(mocks.MockMessageUsecase)
	messageController := controller.NewMessageController(mockMessageUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/chats/:chat_id/messages/:message_id", messageController.DeleteMessage)

	t.Run("success", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		messageID := primitive.NewObjectID()

		mockMessageUsecase.On("DeleteMessage", mock.Anything, chatID, messageID).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/chats/"+chatID.Hex()+"/messages/"+messageID.Hex(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})

	t.Run("invalid IDs", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/chats/invalid-id/messages/invalid-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		messageID := primitive.NewObjectID()

		mockMessageUsecase.On("DeleteMessage", mock.Anything, chatID, messageID).Return(errors.New("deletion error")).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/chats/"+chatID.Hex()+"/messages/"+messageID.Hex(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})
}

func TestUpdateMessage(t *testing.T) {
	mockMessageUsecase := new(mocks.MockMessageUsecase)
	messageController := controller.NewMessageController(mockMessageUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/chats/:chat_id/messages/:message_id", messageController.UpdateMessage)

	t.Run("success", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		messageID := primitive.NewObjectID()
		newContent := "Updated message content"

		mockMessageUsecase.On("UpdateMessage", mock.Anything, chatID, messageID, newContent).Return(nil).Once()

		jsonContent, _ := json.Marshal(newContent)
		req, _ := http.NewRequest(http.MethodPut, "/chats/"+chatID.Hex()+"/messages/"+messageID.Hex(), bytes.NewBuffer(jsonContent))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})

	t.Run("invalid IDs", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/chats/invalid-id/messages/invalid-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		chatID := primitive.NewObjectID()
		messageID := primitive.NewObjectID()

		req, _ := http.NewRequest(http.MethodPut, "/chats/"+chatID.Hex()+"/messages/"+messageID.Hex(), bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

