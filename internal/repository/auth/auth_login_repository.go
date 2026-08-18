package auth

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repositorymodel "github.com/twrnakata/user-auth-api/internal/repository/auth/model"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
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
		return nil, apperror.ErrUserCollectionNil
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
		return apperror.ErrAuthLoginRepositoryNotConfigured
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
