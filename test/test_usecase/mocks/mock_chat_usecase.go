package mocks
import (
	"Real-Time-Chat-Application/domain"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockChatUsecase struct {
	mock.Mock
}

func (m *MockChatUsecase) CreateChat(ctx context.Context, chat *domain.Chat) (primitive.ObjectID, error) {
	args := m.Called(ctx, chat)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockChatUsecase) GetChat(ctx context.Context, chatID primitive.ObjectID) (*domain.Chat, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatUsecase) GetChatsByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Chat, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Chat), args.Error(1)
}

func (m *MockChatUsecase) UpdateChat(ctx context.Context, chatID primitive.ObjectID, chat *domain.Chat) error {
	args := m.Called(ctx, chatID, chat)
	return args.Error(0)
}

func (m *MockChatUsecase) DeleteChat(ctx context.Context, chatID primitive.ObjectID) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}
