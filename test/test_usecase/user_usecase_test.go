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

func TestCreateUser(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	userUsecase := usecase.NewUserUsecase(mockUserRepository, 1*time.Second)

	user := &domain.User{
		UserID:    primitive.NewObjectID(),
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password",
		Chats:     []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUserRepository.On("CreateUser", mock.Anything, user).Return(user.UserID, nil)

	userID, err := userUsecase.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, user.UserID, userID)
	mockUserRepository.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	userUsecase := usecase.NewUserUsecase(mockUserRepository, 1*time.Second)

	userID := primitive.NewObjectID()
	expectedUser := &domain.User{
		UserID:    userID,
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password",
		Chats:     []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUserRepository.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)

	user, err := userUsecase.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepository.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	userUsecase := usecase.NewUserUsecase(mockUserRepository, 1*time.Second)

	email := "test@example.com"
	expectedUser := &domain.User{
		UserID:    primitive.NewObjectID(),
		Email:     email,
		Username:  "testuser",
		Password:  "password",
		Chats:     []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUserRepository.On("GetUserByEmail", mock.Anything, email).Return(expectedUser, nil)

	user, err := userUsecase.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepository.AssertExpectations(t)
}

func TestGetUserByUsername(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	userUsecase := usecase.NewUserUsecase(mockUserRepository, 1*time.Second)

	username := "testuser"
	expectedUser := &domain.User{
		UserID:    primitive.NewObjectID(),
		Email:     "test@example.com",
		Username:  username,
		Password:  "password",
		Chats:     []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUserRepository.On("GetUserByUsername", mock.Anything, username).Return(expectedUser, nil)

	user, err := userUsecase.GetUserByUsername(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserRepository.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	userUsecase := usecase.NewUserUsecase(mockUserRepository, 1*time.Second)

	userID := primitive.NewObjectID()
	updatedUser := &domain.User{
		UserID:    userID,
		Email:     "updated@example.com",
		Username:  "updateduser",
		Password:  "newpassword",
		Chats:     []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUserRepository.On("UpdateUser", mock.Anything, userID, updatedUser).Return(nil)

	err := userUsecase.UpdateUser(context.Background(), userID, updatedUser)
	assert.NoError(t, err)
	mockUserRepository.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	userUsecase := usecase.NewUserUsecase(mockUserRepository, 1*time.Second)

	userID := primitive.NewObjectID()
	mockUserRepository.On("DeleteUser", mock.Anything, userID).Return(nil)

	err := userUsecase.DeleteUser(context.Background(), userID)
	assert.NoError(t, err)
	mockUserRepository.AssertExpectations(t)
}
