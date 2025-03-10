package repository

import (
	"Real-Time-Chat-Application/domain"
	"Real-Time-Chat-Application/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database, collection string) domain.UserRepository {
	return &UserRepository{collection: database.Collection(collection)}
}

func(userrepo *UserRepository) CreateUser(ctx context.Context, user *domain.User) (primitive.ObjectID, error) {

	collection := userrepo.collection

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil

}
func(userrepo *UserRepository) GetUser(ctx context.Context, userID primitive.ObjectID) (*domain.User, error){

	collection := userrepo.collection

	var user domain.User

	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil

}
func(userrepo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error){

	collection := userrepo.collection

	var user domain.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil

}
func(userrepo *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error){

	collection := userrepo.collection

	var user domain.User

	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil

}
func(userrepo *UserRepository) UpdateUser(ctx context.Context, userID primitive.ObjectID, user *domain.User) error{

	collection := userrepo.collection

	updatedFields := bson.M{}

	if user.Username != "" {
		updatedFields["username"] = user.Username
	}

	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("Failed to hash the password %w", err)
		}
		updatedFields["password"] = hashedPassword
	}

	if len(updatedFields) == 0{
		return fmt.Errorf("No fields provided for update")
	}

	update := bson.M{"$set": updatedFields}

	_, err := collection.UpdateOne(ctx,bson.M{"_id":userID},update)

	if err != nil {
		return fmt.Errorf("Failed to update the user %w", err)
	}

	return nil
}
func(userrepo *UserRepository) DeleteUser(ctx context.Context, userID primitive.ObjectID) error{

	collection := userrepo.collection

	_, err := collection.DeleteOne(ctx, bson.M{"_id": userID})

	if err != nil {
		return fmt.Errorf("Failed to delete the user %w", err)
	}

	return nil
}