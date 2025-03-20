package usecase

import (
	"Real-Time-Chat-Application/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository domain.UserRepository, contextTimeout time.Duration) domain.UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
		contextTimeout: contextTimeout,
	}
}

func (userUsecase *UserUsecase) CreateUser(ctx context.Context, user *domain.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, userUsecase.contextTimeout)
	defer cancel()

	userID, err := userUsecase.userRepository.CreateUser(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return userID, nil
}

func (userUsecase *UserUsecase) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, userUsecase.contextTimeout)
	defer cancel()

	user, err := userUsecase.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userUsecase *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, userUsecase.contextTimeout)
	defer cancel()

	user, err := userUsecase.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userUsecase *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, userUsecase.contextTimeout)
	defer cancel()

	user, err := userUsecase.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userUsecase *UserUsecase) UpdateUser(ctx context.Context, userID primitive.ObjectID, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, userUsecase.contextTimeout)
	defer cancel()

	err := userUsecase.userRepository.UpdateUser(ctx, userID, user)
	if err != nil {
		return err
	}
	return nil
}

func (userUsecase *UserUsecase) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, userUsecase.contextTimeout)
	defer cancel()

	err := userUsecase.userRepository.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
