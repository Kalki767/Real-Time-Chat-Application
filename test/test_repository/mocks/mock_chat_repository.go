package mocks
import (
	"Real-Time-Chat-Application/domain"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) CreateChat(ctx context.Context, SenderID primitive.ObjectID, ReceiverID primitive.ObjectID) (primitive.ObjectID, error) {
	args := m.Called(ctx, SenderID, ReceiverID)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockChatRepository) GetChat(ctx context.Context, chatID primitive.ObjectID) (*domain.Chat, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatRepository) GetChatsByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Chat, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Chat), args.Error(1)
}

func (m *MockChatRepository) UpdateChat(ctx context.Context, chatID primitive.ObjectID, chat *domain.Chat) error {
	args := m.Called(ctx, chatID, chat)
	return args.Error(0)
}

func (m *MockChatRepository) DeleteChat(ctx context.Context, chatID primitive.ObjectID) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func (m *MockChatRepository) GetChatByParticipants(ctx context.Context, SenderID primitive.ObjectID, ReceiverID primitive.ObjectID) (*domain.Chat, error) {
	args := m.Called(ctx, SenderID, ReceiverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}
