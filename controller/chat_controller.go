package controller

import (
	"Real-Time-Chat-Application/domain"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatController struct {
	ChatUsecase domain.ChatUsecase
}

func NewChatController(chatUsecase domain.ChatUsecase) *ChatController {
	return &ChatController{ChatUsecase: chatUsecase}
}

func (c *ChatController) CreateChat(ctx *gin.Context) {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var chat domain.Chat
	if err := ctx.ShouldBindJSON(&chat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatID, err := c.ChatUsecase.CreateChat(ctxWithTimeout, &chat)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"chat_id": chatID})
}

func (c *ChatController) GetChat(ctx *gin.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	chatid := ctx.Param("chat_id")

	chatID, err := primitive.ObjectIDFromHex(chatid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	chat, err := c.ChatUsecase.GetChat(ctxWithTimeout, chatID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { // Handle "chat not found" properly
			ctx.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve chat"})
		return
	}

	ctx.JSON(http.StatusOK, chat)
}

func (c *ChatController) GetChatsByUserID(ctx *gin.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	userid := ctx.Param("user_id")

	userID, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	chats, err := c.ChatUsecase.GetChatsByUserID(ctxWithTimeout, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chats)
}

func (c *ChatController) UpdateChat(ctx *gin.Context) {	
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	chatid := ctx.Param("chat_id")

	chatID, err := primitive.ObjectIDFromHex(chatid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var chat domain.Chat
	if err := ctx.ShouldBindJSON(&chat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.ChatUsecase.UpdateChat(ctxWithTimeout, chatID, &chat)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Chat updated successfully"})
}

func (c *ChatController) DeleteChat(ctx *gin.Context) {	
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	chatid := ctx.Param("chat_id")

	chatID, err := primitive.ObjectIDFromHex(chatid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	err = c.ChatUsecase.DeleteChat(ctxWithTimeout, chatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}
