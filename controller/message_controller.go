package controller

import (
	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageController struct {
	messageUsecase domain.MessageUsecase
	hub            *websocket.Hub
}

func NewMessageController(messageUsecase domain.MessageUsecase, hub *websocket.Hub) *MessageController {
	return &MessageController{
		messageUsecase: messageUsecase,
		hub:           hub,
	}
}

// HandleWebSocket handles websocket connections for real-time messaging
func (mc *MessageController) HandleWebSocket(c *gin.Context) {
	websocket.HandleWebSocket(c, mc.hub)
}

// SendMessage handles sending a new message
func (mc *MessageController) SendMessage(c *gin.Context) {
	var message domain.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatID, err := primitive.ObjectIDFromHex(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	err = mc.messageUsecase.SendMessage(c.Request.Context(), chatID, &message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Broadcast the message through websocket
	broadcastMsg := struct {
		ChatID  string         `json:"chat_id"`
		Message domain.Message `json:"message"`
	}{
		ChatID:  chatID.Hex(),
		Message: message,
	}
	mc.hub.BroadcastToChat(chatID.Hex(), []byte(broadcastMsg.Message.Content))

	c.JSON(http.StatusOK, message)
}

// GetMessages retrieves all messages for a chat
func (mc *MessageController) GetMessages(c *gin.Context) {
	chatID, err := primitive.ObjectIDFromHex(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	messages, err := mc.messageUsecase.GetMessages(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetMessage retrieves a specific message
func (mc *MessageController) GetMessage(c *gin.Context) {
	chatID, err := primitive.ObjectIDFromHex(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	messageID, err := primitive.ObjectIDFromHex(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := mc.messageUsecase.GetMessage(c.Request.Context(), chatID, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// DeleteMessage handles message deletion
func (mc *MessageController) DeleteMessage(c *gin.Context) {
	chatID, err := primitive.ObjectIDFromHex(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	messageID, err := primitive.ObjectIDFromHex(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	err = mc.messageUsecase.DeleteMessage(c.Request.Context(), chatID, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// UpdateMessage handles message updates
func (mc *MessageController) UpdateMessage(c *gin.Context) {
	chatID, err := primitive.ObjectIDFromHex(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	messageID, err := primitive.ObjectIDFromHex(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var updateReq struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = mc.messageUsecase.UpdateMessage(c.Request.Context(), chatID, messageID, updateReq.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
}
