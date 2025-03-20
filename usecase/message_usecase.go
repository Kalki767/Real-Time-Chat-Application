package usecase

import (
	"Real-Time-Chat-Application/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageUsecase struct {
	messageRepo domain.MessageRepository
	contextTimeout time.Duration
}

func NewMessageUsecase(messageRepo domain.MessageRepository, contextTimeout time.Duration) domain.MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
		contextTimeout: contextTimeout,
	}
}
func(messageUsecase MessageUsecase) SendMessage(ctx context.Context, chatID primitive.ObjectID, message *domain.Message) error{
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, messageUsecase.contextTimeout)
	defer cancel()

	// Call the repository layer to send the message
	err := messageUsecase.messageRepo.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	return nil

}
func(messageUsecase MessageUsecase) GetMessages(ctx context.Context, chatID primitive.ObjectID) ([]domain.Message, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, messageUsecase.contextTimeout)
	defer cancel()

	// Call the repository layer to get the messages
	messages, err := messageUsecase.messageRepo.GetMessages(ctx,chatID)
	if err != nil{
		return []domain.Message{}, err
	}

	return messages, nil
}
func(messageUsecase MessageUsecase) GetMessage(ctx context.Context, chatID, messageID primitive.ObjectID) (domain.Message, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, messageUsecase.contextTimeout)
	defer cancel()

	// Call the repository layer to get the message
	message, err := messageUsecase.messageRepo.GetMessage(ctx, chatID, messageID)
	if err != nil{
		return domain.Message{}, err
	}

	return message, nil
} 
func(messageUsecase MessageUsecase) DeleteMessage(ctx context.Context, chatID, messageID primitive.ObjectID) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, messageUsecase.contextTimeout)
	defer cancel()

	// Call the repository layer to delete the message
	err := messageUsecase.messageRepo.DeleteMessage(ctx, chatID, messageID)
	if err != nil{
		return err
	}

	return nil

}
func(messageUsecase MessageUsecase) UpdateMessage(ctx context.Context, chatID, messageID primitive.ObjectID, newContent string) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, messageUsecase.contextTimeout)
	defer cancel()

	// Call the repository layer to update the message
	err := messageUsecase.messageRepo.UpdateMessage(ctx, chatID, messageID, newContent)
	if err != nil{
		return err
	}

	return nil

}
