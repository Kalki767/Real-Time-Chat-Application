package usecase

import (
	"Real-Time-Chat-Application/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUsecase struct {
	chatRepository domain.ChatRepository
	contextTimeout time.Duration
}

func NewChatUsecase(chatRepository domain.ChatRepository, timeout time.Duration) domain.ChatUsecase {
	return &ChatUsecase{
		chatRepository: chatRepository,
		contextTimeout: timeout,
	}
}

func (chatusecase *ChatUsecase) CreateChat(ctx context.Context, chat *domain.Chat) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, chatusecase.contextTimeout)
	defer cancel()

	chatID, err := chatusecase.chatRepository.CreateChat(ctx, chat)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return chatID, nil
}

func (chatusecase *ChatUsecase) GetChat(ctx context.Context, chatID primitive.ObjectID) (*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, chatusecase.contextTimeout)
	defer cancel()

	chat, err := chatusecase.chatRepository.GetChat(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (chatusecase *ChatUsecase) GetChatsByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, chatusecase.contextTimeout)
	defer cancel()

	chats, err := chatusecase.chatRepository.GetChatsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (chatusecase *ChatUsecase) UpdateChat(ctx context.Context, chatID primitive.ObjectID, chat *domain.Chat) error {
	ctx, cancel := context.WithTimeout(ctx, chatusecase.contextTimeout)
	defer cancel()

	err := chatusecase.chatRepository.UpdateChat(ctx, chatID, chat)
	if err != nil {
		return err
	}
	return nil
}

func (chatusecase *ChatUsecase) DeleteChat(ctx context.Context, chatID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, chatusecase.contextTimeout)
	defer cancel()

	err := chatusecase.chatRepository.DeleteChat(ctx, chatID)
	if err != nil {
		return err
	}
	return nil
}




