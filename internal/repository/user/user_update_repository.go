package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	repositorymodel "backend-challenge-golang-7solution/internal/repository/user/model"
)

type userDocumentUpdater interface {
	FindOneAndUpdate(executionContext context.Context, filter any, update any, result any) error
}

type mongoUserDocumentUpdater struct {
	userCollection *mongo.Collection
}

func (updater mongoUserDocumentUpdater) FindOneAndUpdate(executionContext context.Context, filter any, update any, result any) error {
	return updater.userCollection.FindOneAndUpdate(
		executionContext,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(result)
}

type UpdateUserRepository struct {
	userCollection      *mongo.Collection
	userDocumentUpdater userDocumentUpdater
}

func NewUpdateUserRepository(userCollection *mongo.Collection) (*UpdateUserRepository, error) {
	if userCollection == nil {
		return nil, errors.New("user collection is nil")
	}

	return &UpdateUserRepository{
		userCollection: userCollection,
		userDocumentUpdater: mongoUserDocumentUpdater{
			userCollection: userCollection,
		},
	}, nil
}

func (repository *UpdateUserRepository) UpdateUser(executionContext context.Context, userID string, update domainuser.UserUpdate, user *domainuser.User) error {
	if repository.userDocumentUpdater == nil {
		return errors.New("update user repository not configured")
	}
	if user == nil {
		return errors.New("user response is nil")
	}

	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidObjectID
	}

	setFields := bson.M{}
	if update.Name != nil {
		setFields["name"] = *update.Name
	}
	if update.Email != nil {
		setFields["email"] = *update.Email
	}

	var document repositorymodel.UpdateUserDocumentModel
	err = repository.userDocumentUpdater.FindOneAndUpdate(
		executionContext,
		bson.M{"_id": objectID},
		bson.M{"$set": setFields},
		&document,
	)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateKey
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
