package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	repositorymodel "github.com/twrnakata/user-auth-api/internal/repository/user/model"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

type userDocumentDeleter interface {
	FindOneAndDelete(executionContext context.Context, filter any, result any) error
}

type mongoUserDocumentDeleter struct {
	userCollection *mongo.Collection
}

func (deleter mongoUserDocumentDeleter) FindOneAndDelete(executionContext context.Context, filter any, result any) error {
	return deleter.userCollection.FindOneAndDelete(executionContext, filter).Decode(result)
}

type DeleteUserRepository struct {
	userCollection      *mongo.Collection
	userDocumentDeleter userDocumentDeleter
}

func NewDeleteUserRepository(userCollection *mongo.Collection) (*DeleteUserRepository, error) {
	if userCollection == nil {
		return nil, apperror.ErrUserCollectionNil
	}

	return &DeleteUserRepository{
		userCollection: userCollection,
		userDocumentDeleter: mongoUserDocumentDeleter{
			userCollection: userCollection,
		},
	}, nil
}

func (repository *DeleteUserRepository) DeleteUser(executionContext context.Context, userID string, user *domainuser.User) error {
	if repository.userDocumentDeleter == nil {
		return apperror.ErrDeleteUserRepositoryNotConfigured
	}
	if user == nil {
		return apperror.ErrUserResponseNil
	}

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidObjectID
	}

	var document repositorymodel.DeleteUserDocumentModel
	err = repository.userDocumentDeleter.FindOneAndDelete(executionContext, bson.M{"_id": objectID}, &document)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		return err
	}

	*user = domainuser.User{
		ID:        document.ID.Hex(),
		Name:      document.Name,
		Email:     document.Email,
		CreatedAt: document.CreatedAt,
	}
	return nil
}
