package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	repositorymodel "backend-challenge-golang-7solution/internal/repository/user/model"
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
		return nil, errors.New("user collection is nil")
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
		return errors.New("delete user repository not configured")
	}
	if user == nil {
		return errors.New("user response is nil")
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
