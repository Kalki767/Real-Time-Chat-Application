// repository/collection.go
package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// CollectionInterface defines the methods used by the repository from mongo.Collection
type CollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions)  SingleResultInterface
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorInterface, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}