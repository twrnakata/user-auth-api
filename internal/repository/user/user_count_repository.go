package user

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

type userDocumentsCounter interface {
	CountDocuments(executionContext context.Context, filter any) (int64, error)
}

type mongoUserDocumentsCounter struct {
	userCollection *mongo.Collection
}

func (counter mongoUserDocumentsCounter) CountDocuments(executionContext context.Context, filter any) (int64, error) {
	return counter.userCollection.CountDocuments(executionContext, filter)
}

type CountUserRepository struct {
	userCollection       *mongo.Collection
	userDocumentsCounter userDocumentsCounter
}

func NewCountUserRepository(userCollection *mongo.Collection) (*CountUserRepository, error) {
	if userCollection == nil {
		return nil, apperror.ErrUserCollectionNil
	}

	return &CountUserRepository{
		userCollection: userCollection,
		userDocumentsCounter: mongoUserDocumentsCounter{
			userCollection: userCollection,
		},
	}, nil
}

func (repository *CountUserRepository) CountUsers(executionContext context.Context, count *int64) error {
	if repository.userDocumentsCounter == nil {
		return apperror.ErrCountUserRepositoryNotConfigured
	}
	if count == nil {
		return apperror.ErrCountResponseNil
	}

	documentsCount, err := repository.userDocumentsCounter.CountDocuments(executionContext, bson.M{})
	if err != nil {
		return err
	}

	*count = documentsCount
	return nil
}
