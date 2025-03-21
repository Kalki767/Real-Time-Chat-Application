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

func TestCreateUser(t *testing.T) {
	mockUserUsecase := new(mocks.MockUserUsecase)
	userController := controller.NewUserController(mockUserUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/users", userController.CreateUser)

	t.Run("success", func(t *testing.T) {
		user := &domain.User{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password123",
		}

		expectedID := primitive.NewObjectID()
		mockUserUsecase.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(expectedID, nil).Once()

		jsonUser, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonUser))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		user := &domain.User{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password123",
		}

		mockUserUsecase.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(primitive.NilObjectID, errors.New("usecase error")).Once()

		jsonUser, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonUser))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestGetUserByID(t *testing.T) {
	mockUserUsecase := new(mocks.MockUserUsecase)
	userController := controller.NewUserController(mockUserUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/users/:id", userController.GetUserByID)

	t.Run("success", func(t *testing.T) {
		userID := primitive.NewObjectID()
		expectedUser := &domain.User{
			UserID:    userID,
			Email:     "test@example.com",
			Username:  "testuser",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserUsecase.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users/"+userID.Hex(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/invalid-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mockUserUsecase.On("GetUserByID", mock.Anything, userID).Return(nil, errors.New("user not found")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users/"+userID.Hex(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestGetUserByEmail(t *testing.T) {
	mockUserUsecase := new(mocks.MockUserUsecase)
	userController := controller.NewUserController(mockUserUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/users/email/:email", userController.GetUserByEmail)

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := &domain.User{
			Email:    email,
			Username: "testuser",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserUsecase.On("GetUserByEmail", mock.Anything, email).Return(expectedUser, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users/email/"+email, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		email := "test@example.com"
		mockUserUsecase.On("GetUserByEmail", mock.Anything, email).Return(nil, errors.New("user not found")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users/email/"+email, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestGetUserByUsername(t *testing.T) {
	mockUserUsecase := new(mocks.MockUserUsecase)
	userController := controller.NewUserController(mockUserUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/users/username/:username", userController.GetUserByUsername)

	t.Run("success", func(t *testing.T) {
		username := "testuser"
		expectedUser := &domain.User{
			Username: username,
			Email:    "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserUsecase.On("GetUserByUsername", mock.Anything, username).Return(expectedUser, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users/username/"+username, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		username := "testuser"
		mockUserUsecase.On("GetUserByUsername", mock.Anything, username).Return(nil, errors.New("user not found")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/users/username/"+username, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockUserUsecase := new(mocks.MockUserUsecase)
	userController := controller.NewUserController(mockUserUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/users/:id", userController.UpdateUser)	

	t.Run("success", func(t *testing.T) {
		userID := primitive.NewObjectID()
		user := &domain.User{
			UserID:    userID,
			Email:     "test@example.com",
			Username:  "testuser",
			Password:  "password123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserUsecase.On("UpdateUser", mock.Anything, userID, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == user.Email && 
				   u.Username == user.Username &&
				   u.UserID == userID
		})).Return(nil).Once()

		jsonUser, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPut, "/users/"+userID.Hex(), bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/users/invalid-id", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		user := &domain.User{
			UserID:    userID,
			Email:     "test@example.com", 
			Username:  "testuser",
			Password:  "password123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserUsecase.On("UpdateUser", mock.Anything, userID, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == user.Email && 
				   u.Username == user.Username &&
				   u.UserID == userID
		})).Return(errors.New("usecase error")).Once()

		jsonUser, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPut, "/users/"+userID.Hex(), bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	mockUserUsecase := new(mocks.MockUserUsecase)
	userController := controller.NewUserController(mockUserUsecase)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/users/:id", userController.DeleteUser)

	t.Run("success", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mockUserUsecase.On("DeleteUser", mock.Anything, userID).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/users/"+userID.Hex(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/invalid-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mockUserUsecase.On("DeleteUser", mock.Anything, userID).Return(errors.New("usecase error")).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/users/"+userID.Hex(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}


	
	
	
	
	



