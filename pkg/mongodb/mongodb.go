package mongodb

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const UsersCollectionName = "users"

func Connect(executionContext context.Context, uri string) (*mongo.Client, error) {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return nil, errors.New("mongo uri is empty")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(executionContext, nil); err != nil {
		_ = client.Disconnect(executionContext)
		return nil, err
	}

	return client, nil
}

func UserCollection(client *mongo.Client, databaseName string) *mongo.Collection {
	return client.Database(databaseName).Collection(UsersCollectionName)
}
