package auth

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
)

var ErrNotFound = errors.New("not found")

type userDocumentFinder interface {
	FindOne(executionContext context.Context, filter any, result any) error
}

type mongoUserDocumentFinder struct {
	userCollection *mongo.Collection
}

func (finder mongoUserDocumentFinder) FindOne(executionContext context.Context, filter any, result any) error {
	return finder.userCollection.FindOne(executionContext, filter).Decode(result)
}

type AuthLoginRepository struct {
	userCollection     *mongo.Collection
	userDocumentFinder userDocumentFinder
}

func NewAuthLoginRepository(userCollection *mongo.Collection) (*AuthLoginRepository, error) {
	if userCollection == nil {
		return nil, errors.New("user collection is nil")
	}

	return &AuthLoginRepository{
		userCollection: userCollection,
		userDocumentFinder: mongoUserDocumentFinder{
			userCollection: userCollection,
		},
	}, nil
}

func (repository *AuthLoginRepository) GetLoginUserByEmail(executionContext context.Context, request *repositorymodel.AuthLoginRepositoryRequestModel, response *repositorymodel.GetLoginUserByEmailModel) error {
	if repository.userDocumentFinder == nil {
		return errors.New("auth login repository not configured")
	}

	var document repositorymodel.GetLoginUserDocumentModel
	err := repository.userDocumentFinder.FindOne(executionContext, bson.M{"email": request.Email}, &document)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound
		}
		return err
	}

	response.ID = document.ID.Hex()
	response.Name = document.Name
	response.Email = document.Email
	response.PasswordHash = document.PasswordHash
	return nil
}
