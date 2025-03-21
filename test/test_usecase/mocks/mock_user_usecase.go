package mocks

import (
	"Real-Time-Chat-Application/domain"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) CreateUser(ctx context.Context, user *domain.User) (primitive.ObjectID, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockUserUsecase) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, userID primitive.ObjectID, user *domain.User) error {
	args := m.Called(ctx, userID, user)
	return args.Error(0)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
