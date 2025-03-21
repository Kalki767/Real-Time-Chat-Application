package mocks

import (
	"Real-Time-Chat-Application/domain"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockMessageUsecase struct {
	mock.Mock
}

func (m *MockMessageUsecase) SendMessage(ctx context.Context, chatID primitive.ObjectID, message *domain.Message) error {
	args := m.Called(ctx, chatID, message)
	return args.Error(0)
}

func (m *MockMessageUsecase) GetMessage(ctx context.Context, chatID, messageID primitive.ObjectID) (domain.Message, error) {
	args := m.Called(ctx, chatID, messageID)
	if args.Get(0) == nil {
		return domain.Message{}, args.Error(1)
	}
	return args.Get(0).(domain.Message), args.Error(1)
}

func (m *MockMessageUsecase) GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]domain.Message, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageUsecase) UpdateMessage(ctx context.Context, chatID, messageID primitive.ObjectID, newContent string) error {
	args := m.Called(ctx, chatID, messageID, newContent)
	return args.Error(0)
}

func (m *MockMessageUsecase) DeleteMessage(ctx context.Context, chatID, messageID primitive.ObjectID) error {
	args := m.Called(ctx, chatID, messageID)
	return args.Error(0)
}
