package test

import (
	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/repository"
	"Real-Time-Chat-Application/test/mongo/mocks"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestSendMessage(t *testing.T) {
	// Setup
	mockCollection := new(mocks.MockCollection)
	repo := repository.NewMessageRepository(mockCollection) 

	chatID := primitive.NewObjectID()
	message := &domain.Message{
		MessageID: primitive.NewObjectID(),
		Time:      time.Now(),
	}

	// Mock UpdateOne result
	updateResult := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}

	// Use mock.MatchedBy to compare bson.M dynamically
	mockCollection.On("UpdateOne", mock.Anything, bson.M{"_id": chatID}, mock.MatchedBy(func(update bson.M) bool {
		// Ensure that "messages" is being pushed and "updated_at" is set
		_, msgExists := update["$push"].(bson.M)["messages"]
		_, timeExists := update["$set"].(bson.M)["updated_at"]
		return msgExists && timeExists
	})).Return(updateResult, nil)

	// Execute
	err := repo.SendMessage(context.TODO(), chatID, message)

	// Verify
	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
}

func TestGetMessage(t *testing.T) {
	// Create a mock collection and single result
	mockCollection := new(mocks.MockCollection)
	mockSingleResult := new(mocks.MockSingleResult)
	chatRepo := repository.NewMessageRepository(mockCollection)

	// Generate common ObjectIDs for consistency
	commonChatID := primitive.NewObjectID()
	commonMessageID := primitive.NewObjectID()

	// Define test cases
	tests := []struct {
		name         string
		chatID       primitive.ObjectID
		messageID    primitive.ObjectID
		mockReturn   *domain.Message
		mockError    error
		expectedMsg  domain.Message
		expectedErr  error
	}{
		{
			name:       "Successful retrieval",
			chatID:     commonChatID, // Use common ID
			messageID:  commonMessageID,
			mockReturn: &domain.Message{MessageID: commonMessageID, Content: "Hello"},
			mockError:  nil,
			expectedMsg: domain.Message{MessageID: commonMessageID, Content: "Hello"}, // Match with mockReturn
			expectedErr: nil,
		},
		{
			name:       "Message not found",
			chatID:     primitive.NewObjectID(),
			messageID:  primitive.NewObjectID(),
			mockReturn: nil,
			mockError:  mongo.ErrNoDocuments,
			expectedMsg: domain.Message{},
			expectedErr: fmt.Errorf("failed to fetch chat: %w", mongo.ErrNoDocuments),
		},
		{
			name:       "Database error",
			chatID:     primitive.NewObjectID(),
			messageID:  primitive.NewObjectID(),
			mockReturn: nil,
			mockError:  errors.New("some database error"),
			expectedMsg: domain.Message{},
			expectedErr: fmt.Errorf("failed to fetch chat: some database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks before each test
			mockCollection = new(mocks.MockCollection)
			mockSingleResult = new(mocks.MockSingleResult)
			chatRepo = repository.NewMessageRepository(mockCollection)

			// Mock the expected FindOne call
			mockCollection.On("FindOne", mock.Anything, bson.M{"_id": tt.chatID, "messages.messageID": tt.messageID}, mock.Anything).
				Return(mockSingleResult)

			// If message exists, mock the Decode method to return it
			if tt.mockReturn != nil {
				mockSingleResult.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]domain.Message)
					*arg = []domain.Message{*tt.mockReturn} // Store the mock message in the slice
				})
			} else {
				mockSingleResult.On("Decode", mock.Anything).Return(tt.mockError)
			}

			// Execute function
			msg, err := chatRepo.GetMessage(context.Background(), tt.chatID, tt.messageID)

			// Assertions
			assert.Equal(t, tt.expectedMsg, msg)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify expectations
			mockCollection.AssertExpectations(t)
			mockSingleResult.AssertExpectations(t)
		})
	}
}
func TestGetMessages(t *testing.T) {
	var (
		mockCollection   *mocks.MockCollection
		mockSingleResult *mocks.MockSingleResult
		
	)

	tests := []struct {
		name        string
		chatID      primitive.ObjectID
		mockReturn  *domain.Chat
		mockError   error
		expectedMsg []domain.Message
		expectedErr error
	}{
		{
			name:   "Success - Messages found",
			chatID: primitive.NewObjectID(),
			mockReturn: &domain.Chat{
				Messages: []domain.Message{
					{
						MessageID: primitive.NewObjectID(),
						Content:   "Message 1",
						Time:     time.Now().Add(-time.Hour),
					},
					{
						MessageID: primitive.NewObjectID(),
						Content:   "Message 2", 
						Time:     time.Now(),
					},
				},
			},
			mockError: nil,
			expectedMsg: []domain.Message{
				{
					MessageID: primitive.NewObjectID(),
					Content:   "Message 1",
					Time:     time.Now().Add(-time.Hour),
				},
				{
					MessageID: primitive.NewObjectID(),
					Content:   "Message 2",
					Time:     time.Now(), 
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Database error",
			chatID:      primitive.NewObjectID(),
			mockReturn:  nil,
			mockError:   errors.New("database error"),
			expectedMsg: []domain.Message{},
			expectedErr: fmt.Errorf("failed to fetch chat: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks before each test
			mockCollection = new(mocks.MockCollection)
			mockSingleResult = new(mocks.MockSingleResult)
			chatRepo := repository.NewMessageRepository(mockCollection)

			// Mock FindOne
			mockCollection.On("FindOne", mock.Anything, bson.M{"_id": tt.chatID}).Return(mockSingleResult)

			if tt.mockReturn != nil {
				mockSingleResult.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*domain.Chat)
					*arg = *tt.mockReturn
				})
			} else {
				mockSingleResult.On("Decode", mock.Anything).Return(tt.mockError)
			}

			// Execute
			messages, err := chatRepo.GetMessages(context.Background(), tt.chatID)

			// Verify
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedMsg), len(messages))
			}

			mockCollection.AssertExpectations(t)
			mockSingleResult.AssertExpectations(t)
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	var (
		mockCollection *mocks.MockCollection
	)

	tests := []struct {
		name        string
		chatID      primitive.ObjectID
		messageID   primitive.ObjectID
		mockResult  *mongo.UpdateResult
		mockError   error
		expectedErr error
	}{
		{
			name:      "Success - Message deleted",
			chatID:    primitive.NewObjectID(),
			messageID: primitive.NewObjectID(),
			mockResult: &mongo.UpdateResult{
				ModifiedCount: 1,
			},
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:      "Failure - Message not found",
			chatID:    primitive.NewObjectID(),
			messageID: primitive.NewObjectID(),
			mockResult: &mongo.UpdateResult{
				ModifiedCount: 0,
			},
			mockError:   nil,
			expectedErr: fmt.Errorf("message not found or already deleted"),
		},
		{
			name:       "Database error",
			chatID:     primitive.NewObjectID(),
			messageID:  primitive.NewObjectID(),
			mockResult: nil,
			mockError:  errors.New("database error"),
			expectedErr: fmt.Errorf("failed to delete message: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockCollection = new(mocks.MockCollection)
			chatRepo := repository.NewMessageRepository(mockCollection)

			// Mock UpdateOne with MatchedBy to handle time.Now()
			mockCollection.On("UpdateOne", 
				mock.Anything,
				bson.M{"_id": tt.chatID},
				mock.MatchedBy(func(update bson.M) bool {
					pull := update["$pull"].(bson.M)
					messages := pull["messages"].(bson.M)
					return messages["message_id"] == tt.messageID && update["$set"] != nil
				}),
			).Return(tt.mockResult, tt.mockError)

			// Execute
			err := chatRepo.DeleteMessage(context.Background(), tt.chatID, tt.messageID)

			// Verify
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockCollection.AssertExpectations(t)
		})
	}
}

func TestUpdateMessage(t *testing.T) {
	var (
		mockCollection *mocks.MockCollection
	)

	tests := []struct {
		name        string
		chatID      primitive.ObjectID
		messageID   primitive.ObjectID
		content     string
		mockResult  *mongo.UpdateResult
		mockError   error
		expectedErr error
	}{
		{
			name:      "Success - Message updated",
			chatID:    primitive.NewObjectID(),
			messageID: primitive.NewObjectID(),
			content:   "Updated content",
			mockResult: &mongo.UpdateResult{
				ModifiedCount: 1,
			},
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name:      "Failure - Message not found",
			chatID:    primitive.NewObjectID(),
			messageID: primitive.NewObjectID(),
			content:   "Updated content",
			mockResult: &mongo.UpdateResult{
				ModifiedCount: 0,
			},
			mockError:   nil,
			expectedErr: fmt.Errorf("message not found or already updated"),
		},
		{
			name:      "Database error",
			chatID:    primitive.NewObjectID(),
			messageID: primitive.NewObjectID(),
			content:   "Updated content",
			mockResult:  nil,
			mockError:   errors.New("database error"),
			expectedErr: fmt.Errorf("failed to update message: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockCollection = new(mocks.MockCollection)
			chatRepo := repository.NewMessageRepository(mockCollection)

			// Mock UpdateOne correctly with 4 parameters
			mockCollection.On("UpdateOne",
				mock.Anything, // Context
				bson.M{"_id": tt.chatID}, // Filter
				mock.MatchedBy(func(update bson.M) bool {
					set, ok := update["$set"].(bson.M)
					if !ok {
						return false
					}
					delete(set, "updated_at") // Ignore the timestamp comparison
					return true
				}),
				mock.Anything, // Options (like arrayFilters)
			).Return(tt.mockResult, tt.mockError)

			// Execute the function being tested
			err := chatRepo.UpdateMessage(context.Background(), tt.chatID, tt.messageID, tt.content)

			// Verify the expected result
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			// Ensure the mock expectations are met
			mockCollection.AssertExpectations(t)
		})
	}
}