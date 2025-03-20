package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/repository"
	"Real-Time-Chat-Application/test/mocks"
)

func TestCreateUser(t *testing.T) {
	t.Run("Success - Creates New User", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		repo := repository.NewUserRepository(mockCollection)

		user := &domain.User{
			Username:  "testuser",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		insertResult := &mongo.InsertOneResult{
			InsertedID: primitive.NewObjectID(),
		}

		// Mock InsertOne
		mockCollection.On("InsertOne", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Username == user.Username && u.Password == user.Password
		})).Return(insertResult, nil)

		// Execute
		userID, err := repo.CreateUser(context.TODO(), user)
		// Verify returned userID matches the inserted ID
		assert.Equal(t, insertResult.InsertedID.(primitive.ObjectID), userID)

		// Verify
		assert.NoError(t, err)
		mockCollection.AssertExpectations(t)
	})

	t.Run("Failure - Duplicate Username", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		repo := repository.NewUserRepository(mockCollection)

		user := &domain.User{
			Username:  "existinguser", 
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock InsertOne returning duplicate key error
		mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{}, mongo.WriteException{
			WriteErrors: []mongo.WriteError{
				{
					Code:    11000,
					Message: "write exception: write errors: [duplicate username]",
				},
			},
		})

		// Execute
		_, err := repo.CreateUser(context.TODO(), user)

		// Verify
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "write exception: write errors: [duplicate username]")
		mockCollection.AssertExpectations(t)
	})
}

func TestGetUserByUsername(t *testing.T) {
	t.Run("Success - User Found", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		mockSingleResult := new(mocks.MockSingleResult)
		repo := repository.NewUserRepository(mockCollection)

		expectedUser := &domain.User{
			UserID:    primitive.NewObjectID(),
			Username:  "testuser",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock Decode
		mockSingleResult.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*domain.User)
			*arg = *expectedUser
		})

		// Mock FindOne
		mockCollection.On("FindOne", mock.Anything, bson.M{"username": expectedUser.Username}).Return(mockSingleResult)

		// Execute
		user, err := repo.GetUserByUsername(context.TODO(), expectedUser.Username)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.Password, user.Password)
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})

	t.Run("Failure - User Not Found", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		mockSingleResult := new(mocks.MockSingleResult)
		repo := repository.NewUserRepository(mockCollection)

		username := "nonexistentuser"

		// Mock FindOne returning an error
		mockSingleResult.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments)
		mockCollection.On("FindOne", mock.Anything, bson.M{"username": username}).Return(mockSingleResult)

		// Execute
		user, err := repo.GetUserByUsername(context.TODO(), username)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, mongo.ErrNoDocuments.Error(), err.Error())
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("Success - User Found", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		mockSingleResult := new(mocks.MockSingleResult)
		repo := repository.NewUserRepository(mockCollection)

		userID := primitive.NewObjectID()
		expectedUser := &domain.User{
			UserID:    userID,
			Username:  "testuser",
			Password:  "hashedpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock Decode
		mockSingleResult.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*domain.User)
			*arg = *expectedUser
		})

		// Mock FindOne
		mockCollection.On("FindOne", mock.Anything, bson.M{"_id": userID}).Return(mockSingleResult)

		// Execute
		user, err := repo.GetUserByID(context.TODO(), userID)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.UserID, user.UserID)
		assert.Equal(t, expectedUser.Username, user.Username)
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})

	t.Run("Failure - User Not Found", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		mockSingleResult := new(mocks.MockSingleResult)
		repo := repository.NewUserRepository(mockCollection)

		userID := primitive.NewObjectID()

		// Mock FindOne returning an error
		mockSingleResult.On("Decode", mock.Anything).Return(mongo.ErrNoDocuments)
		mockCollection.On("FindOne", mock.Anything, bson.M{"_id": userID}).Return(mockSingleResult)

		// Execute
		user, err := repo.GetUserByID(context.TODO(), userID)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, mongo.ErrNoDocuments.Error(), err.Error())
		mockCollection.AssertExpectations(t)
		mockSingleResult.AssertExpectations(t)
	})
}
func TestUpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		repo := repository.NewUserRepository(mockCollection)

		userID := primitive.NewObjectID()
		updateUser := &domain.User{
			Username: "newusername",
			Password: "newpassword",
		}


		// Mock UpdateOne
		mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": userID}, mock.MatchedBy(func(update bson.M) bool {
			// Verify the update has the correct structure
			setMap := update["$set"].(bson.M)
			return setMap["username"] == "newusername" && setMap["password"] != ""
		})).Return(&mongo.UpdateResult{ModifiedCount: 1}, nil)

		// Execute
		err := repo.UpdateUser(context.TODO(), userID, updateUser)

		// Verify
		assert.NoError(t, err)
		mockCollection.AssertExpectations(t)
	})

	t.Run("Failure - No Fields to Update", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		repo := repository.NewUserRepository(mockCollection)

		userID := primitive.NewObjectID()
		updateUser := &domain.User{} // Empty user

		// Execute
		err := repo.UpdateUser(context.TODO(), userID, updateUser)

		// Verify
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "No fields provided for update")
		mockCollection.AssertExpectations(t)
	})

	t.Run("Failure - Update Error", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		repo := repository.NewUserRepository(mockCollection)

		userID := primitive.NewObjectID()
		updateUser := &domain.User{
			Username: "newusername",
		}

		// Mock UpdateOne returning an error
		mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": userID}, mock.Anything).
			Return(&mongo.UpdateResult{}, fmt.Errorf("update error"))

		// Execute
		err := repo.UpdateUser(context.TODO(), userID, updateUser)

		// Verify
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to update the user")
		mockCollection.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		repo := repository.NewUserRepository(mockCollection)

		userID := primitive.NewObjectID()

		// Mock DeleteOne
		mockCollection.On("DeleteOne", mock.Anything, bson.M{"_id": userID}).
			Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

		// Execute
		err := repo.DeleteUser(context.TODO(), userID)

		// Verify
		assert.NoError(t, err)
		mockCollection.AssertExpectations(t)
	})

	t.Run("Failure - Delete Error", func(t *testing.T) {
		// Setup
		mockCollection := new(mocks.MockCollection)
		repo := repository.NewUserRepository(mockCollection)

		userID := primitive.NewObjectID()

		// Mock DeleteOne returning an error
		mockCollection.On("DeleteOne", mock.Anything, bson.M{"_id": userID}).
			Return(&mongo.DeleteResult{}, fmt.Errorf("delete error"))

		// Execute
		err := repo.DeleteUser(context.TODO(), userID)

		// Verify
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to delete the user")
		mockCollection.AssertExpectations(t)
	})
}
