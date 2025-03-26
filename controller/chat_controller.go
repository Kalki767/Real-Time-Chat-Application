package controller
import (
	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatController struct {
	chatUsecase domain.ChatUsecase
	hub         *websocket.Hub
}

func NewChatController(chatUsecase domain.ChatUsecase, hub *websocket.Hub) *ChatController {
	return &ChatController{
		chatUsecase: chatUsecase,
		hub:         hub,
	}
}

// CreateChat handles creation of a new chat
func (cc *ChatController) CreateChat(c *gin.Context) {
	var chatRequest struct {
		SenderID   primitive.ObjectID `json:"sender_id"`
		ReceiverID primitive.ObjectID `json:"receiver_id"`
	}
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatID, err := cc.chatUsecase.CreateChat(c.Request.Context(), chatRequest.SenderID, chatRequest.ReceiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat_id": chatID})
}

// GetChat retrieves a specific chat
func (cc *ChatController) GetChat(c *gin.Context) {
	chatID, err := primitive.ObjectIDFromHex(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	chat, err := cc.chatUsecase.GetChat(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// GetUserChats retrieves all chats for a user
func (cc *ChatController) GetUserChats(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	chats, err := cc.chatUsecase.GetChatsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// DeleteChat handles chat deletion
func (cc *ChatController) DeleteChat(c *gin.Context) {
	chatID, err := primitive.ObjectIDFromHex(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	err = cc.chatUsecase.DeleteChat(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}

// UpdateChat handles chat updates
func (cc *ChatController) UpdateChat(c *gin.Context) {
	chatID, err := primitive.ObjectIDFromHex(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var chat domain.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cc.chatUsecase.UpdateChat(c.Request.Context(), chatID, &chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (cc *ChatController) GetChatByParticipants(c *gin.Context) {
	SenderID, err := primitive.ObjectIDFromHex(c.Param("sender_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
		return
	}

	ReceiverID, err := primitive.ObjectIDFromHex(c.Param("receiver_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
		return
	}

	chat, err := cc.chatUsecase.GetChatByParticipants(c.Request.Context(), SenderID, ReceiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}


