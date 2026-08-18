package mongodb

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

const UsersCollectionName = "users"

func Connect(executionContext context.Context, uri string) (*mongo.Client, error) {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return nil, apperror.ErrMongoURIEmpty
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
