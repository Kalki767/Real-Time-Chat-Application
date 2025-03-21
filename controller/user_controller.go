package controller

import (
	"Real-Time-Chat-Application/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	UserUsecase domain.UserUsecase
}

func NewUserController(us domain.UserUsecase) *UserController {
	return &UserController{
		UserUsecase: us,
	}
}

func (c *UserController) CreateUser(context *gin.Context) {
	var user domain.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	userID, err := c.UserUsecase.CreateUser(context.Request.Context(), &user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	context.JSON(http.StatusCreated, gin.H{"user_id": userID})
	
}

func (c *UserController) GetUserByID(context *gin.Context) {
	params := context.Param("id")

	userID, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.UserUsecase.GetUserByID(context.Request.Context(), userID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}
func (c *UserController) GetUserByEmail(context *gin.Context) {
	email := context.Param("email")

	user, err := c.UserUsecase.GetUserByEmail(context.Request.Context(), email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}

func (c *UserController) GetUserByUsername(context *gin.Context) {
	username := context.Param("username")

	user, err := c.UserUsecase.GetUserByUsername(context.Request.Context(), username)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}


func (c *UserController) UpdateUser(context *gin.Context) {
	params := context.Param("id")
	userID, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user domain.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user.UserID = userID
	user.UpdatedAt = time.Now()

	err = c.UserUsecase.UpdateUser(context.Request.Context(), userID, &user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (c *UserController) DeleteUser(context *gin.Context) {
	params := context.Param("id")
	userID, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = c.UserUsecase.DeleteUser(context.Request.Context(), userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
