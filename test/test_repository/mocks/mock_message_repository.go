package mocks

import (
	"context"
	"Real-Time-Chat-Application/domain"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) SendMessage(ctx context.Context, chatID primitive.ObjectID, message *domain.Message) error {
	args := m.Called(ctx, chatID, message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]domain.Message, error) {
	args := m.Called(ctx, chatID)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetMessage(ctx context.Context, chatID, messageID primitive.ObjectID) (domain.Message, error) {
	args := m.Called(ctx, chatID, messageID)
	return args.Get(0).(domain.Message), args.Error(1)
}

func (m *MockMessageRepository) DeleteMessage(ctx context.Context, chatID, messageID primitive.ObjectID) error {
	args := m.Called(ctx, chatID, messageID)
	return args.Error(0)
}

func (m *MockMessageRepository) UpdateMessage(ctx context.Context, chatID, messageID primitive.ObjectID, newContent string) error {
	args := m.Called(ctx, chatID, messageID, newContent)
	return args.Error(0)
}
