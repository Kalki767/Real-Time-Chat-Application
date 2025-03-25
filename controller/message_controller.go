package controller

import (
	"Real-Time-Chat-Application/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type MessageController struct {
	MessageUsecase domain.MessageUsecase
}

func NewMessageController(messageUsecase domain.MessageUsecase) *MessageController {
	return &MessageController{MessageUsecase: messageUsecase}
}

func (c *MessageController) SendMessage(ctx *gin.Context) {
	params := ctx.Param("chat_id")

	chatID, err := primitive.ObjectIDFromHex(params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var message domain.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.MessageUsecase.SendMessage(ctx.Request.Context(), chatID, &message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Message sent successfully"})
}

func (c *MessageController) GetMessages(ctx *gin.Context) {
	params := ctx.Param("chat_id")

	chatID, err := primitive.ObjectIDFromHex(params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	messages,err := c.MessageUsecase.GetMessages(ctx.Request.Context(), chatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"messages": messages})
	
}

func (c *MessageController) GetMessage(ctx *gin.Context) {
	chatIDParam := ctx.Param("chat_id")
	messageIDParam := ctx.Param("message_id")

	chatID, err := primitive.ObjectIDFromHex(chatIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	messageID, err := primitive.ObjectIDFromHex(messageIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := c.MessageUsecase.GetMessage(ctx.Request.Context(), chatID, messageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func (c *MessageController) DeleteMessage(ctx *gin.Context) {
	chatIDParam := ctx.Param("chat_id")
	messageIDParam := ctx.Param("message_id")
	
	chatID, err := primitive.ObjectIDFromHex(chatIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}	

	messageID, err := primitive.ObjectIDFromHex(messageIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}	

	err = c.MessageUsecase.DeleteMessage(ctx.Request.Context(), chatID, messageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

func (c *MessageController) UpdateMessage(ctx *gin.Context) {
	chatIDParam := ctx.Param("chat_id")
	messageIDParam := ctx.Param("message_id")
	
	chatID, err := primitive.ObjectIDFromHex(chatIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}		

	messageID, err := primitive.ObjectIDFromHex(messageIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}		

	var newContent string
	if err := ctx.ShouldBindJSON(&newContent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	

	err = c.MessageUsecase.UpdateMessage(ctx.Request.Context(), chatID, messageID, newContent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	
}

