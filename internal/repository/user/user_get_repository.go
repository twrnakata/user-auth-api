package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repositorymodel "backend-challenge-golang-7solution/internal/repository/user/model"
	servicemodel "backend-challenge-golang-7solution/internal/service/user/model"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidObjectID = errors.New("invalid object id")
)

type userDocumentFinder interface {
	FindOne(executionContext context.Context, filter any, result any) error
}

type mongoUserDocumentFinder struct {
	userCollection *mongo.Collection
}

func (finder mongoUserDocumentFinder) FindOne(executionContext context.Context, filter any, result any) error {
	return finder.userCollection.FindOne(executionContext, filter).Decode(result)
}

type GetUserRepository struct {
	userCollection     *mongo.Collection
	userDocumentFinder userDocumentFinder
}

func NewGetUserRepository(userCollection *mongo.Collection) (*GetUserRepository, error) {
	if userCollection == nil {
		return nil, errors.New("user collection is nil")
	}

	return &GetUserRepository{
		userCollection: userCollection,
		userDocumentFinder: mongoUserDocumentFinder{
			userCollection: userCollection,
		},
	}, nil
}

func (repository *GetUserRepository) GetUserByID(executionContext context.Context, request *servicemodel.GetUserRequestModel, response *repositorymodel.GetUserByIDModel) error {
	if repository.userDocumentFinder == nil {
		return errors.New("get user repository not configured")
	}

	objectID, err := bson.ObjectIDFromHex(request.ID)
	if err != nil {
		return ErrInvalidObjectID
	}

	var document repositorymodel.GetUserByIDDocumentModel
	err = repository.userDocumentFinder.FindOne(executionContext, bson.M{"_id": objectID}, &document)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		return err
	}

	response.ID = document.ID.Hex()
	response.Name = document.Name
	response.Email = document.Email
	response.CreatedAt = document.CreatedAt
	return nil
}
