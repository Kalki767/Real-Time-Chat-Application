package test
import (
	"Real-Time-Chat-Application/controller"
	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/test/test_usecase/mocks"
	"Real-Time-Chat-Application/websocket"
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

func TestSendMessage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockMessageUsecase := new(mocks.MockMessageUsecase)
		hub := websocket.NewHub()
		messageController := controller.NewMessageController(mockMessageUsecase, hub)

		r := gin.Default()
		r.POST("/chats/:chat_id/messages", messageController.SendMessage)

		chatID := primitive.NewObjectID()
		message := domain.Message{
			Content:  "Test message",
			SenderID: primitive.NewObjectID(),
		}

		mockMessageUsecase.On("SendMessage", mock.Anything, chatID, mock.AnythingOfType("*domain.Message")).Return(nil)

		jsonValue, _ := json.Marshal(message)
		req, _ := http.NewRequest("POST", "/chats/"+chatID.Hex()+"/messages", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})

	t.Run("Invalid Chat ID", func(t *testing.T) {
		mockMessageUsecase := new(mocks.MockMessageUsecase)
		hub := websocket.NewHub()
		messageController := controller.NewMessageController(mockMessageUsecase, hub)

		r := gin.Default()
		r.POST("/chats/:chat_id/messages", messageController.SendMessage)

		message := domain.Message{
			Content:  "Test message",
			SenderID: primitive.NewObjectID(),
		}

		jsonValue, _ := json.Marshal(message)
		req, _ := http.NewRequest("POST", "/chats/invalid_id/messages", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockMessageUsecase := new(mocks.MockMessageUsecase)
		hub := websocket.NewHub()
		messageController := controller.NewMessageController(mockMessageUsecase, hub)

		r := gin.Default()
		r.POST("/chats/:chat_id/messages", messageController.SendMessage)

		chatID := primitive.NewObjectID()
		message := domain.Message{
			Content:  "Test message",
			SenderID: primitive.NewObjectID(),
		}

		mockMessageUsecase.On("SendMessage", mock.Anything, chatID, mock.AnythingOfType("*domain.Message")).Return(errors.New("internal error"))

		jsonValue, _ := json.Marshal(message)
		req, _ := http.NewRequest("POST", "/chats/"+chatID.Hex()+"/messages", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})
}

func TestGetMessages(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockMessageUsecase := new(mocks.MockMessageUsecase)
		hub := websocket.NewHub()
		messageController := controller.NewMessageController(mockMessageUsecase, hub)

		r := gin.Default()
		r.GET("/chats/:chat_id/messages", messageController.GetMessages)

		chatID := primitive.NewObjectID()
		expectedMessages := []domain.Message{
			{
				MessageID: primitive.NewObjectID(),
				Content:   "Test message 1",
				SenderID:  primitive.NewObjectID(),
			},
			{
				MessageID: primitive.NewObjectID(),
				Content:   "Test message 2", 
				SenderID:  primitive.NewObjectID(),
			},
		}

		mockMessageUsecase.On("GetMessages", mock.Anything, chatID).Return(expectedMessages, nil)

		req, _ := http.NewRequest("GET", "/chats/"+chatID.Hex()+"/messages", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []domain.Message
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedMessages), len(response))
		mockMessageUsecase.AssertExpectations(t)
	})

	t.Run("Invalid Chat ID", func(t *testing.T) {
		mockMessageUsecase := new(mocks.MockMessageUsecase)
		hub := websocket.NewHub()
		messageController := controller.NewMessageController(mockMessageUsecase, hub)

		r := gin.Default()
		r.GET("/chats/:chat_id/messages", messageController.GetMessages)

		req, _ := http.NewRequest("GET", "/chats/invalid_id/messages", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteMessage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockMessageUsecase := new(mocks.MockMessageUsecase)
		hub := websocket.NewHub()
		messageController := controller.NewMessageController(mockMessageUsecase, hub)

		r := gin.Default()
		r.DELETE("/chats/:chat_id/messages/:message_id", messageController.DeleteMessage)

		chatID := primitive.NewObjectID()
		messageID := primitive.NewObjectID()

		mockMessageUsecase.On("DeleteMessage", mock.Anything, chatID, messageID).Return(nil)

		req, _ := http.NewRequest("DELETE", "/chats/"+chatID.Hex()+"/messages/"+messageID.Hex(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})
}

func TestUpdateMessage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockMessageUsecase := new(mocks.MockMessageUsecase)
		hub := websocket.NewHub()
		messageController := controller.NewMessageController(mockMessageUsecase, hub)

		r := gin.Default()
		r.PUT("/chats/:chat_id/messages/:message_id", messageController.UpdateMessage)

		chatID := primitive.NewObjectID()
		messageID := primitive.NewObjectID()
		updateReq := struct {
			Content string `json:"content"`
		}{
			Content: "Updated message",
		}

		mockMessageUsecase.On("UpdateMessage", mock.Anything, chatID, messageID, updateReq.Content).Return(nil)

		jsonValue, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest("PUT", "/chats/"+chatID.Hex()+"/messages/"+messageID.Hex(), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockMessageUsecase.AssertExpectations(t)
	})
}
