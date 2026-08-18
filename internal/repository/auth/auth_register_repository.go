package auth

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	repositorymodel "github.com/twrnakata/user-auth-api/internal/repository/auth/model"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
	"github.com/twrnakata/user-auth-api/pkg/datetime"
)

var ErrDuplicateKey = errors.New("duplicate key")

type userDocumentInserter interface {
	InsertOne(executionContext context.Context, document any) (insertedID any, err error)
}

type mongoUserDocumentInserter struct {
	userCollection *mongo.Collection
}

func (inserter mongoUserDocumentInserter) InsertOne(executionContext context.Context, document any) (any, error) {
	result, err := inserter.userCollection.InsertOne(executionContext, document)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

type AuthRegisterRepository struct {
	userCollection       *mongo.Collection
	userDocumentInserter userDocumentInserter
}

func NewAuthRegisterRepository(executionContext context.Context, userCollection *mongo.Collection) (*AuthRegisterRepository, error) {
	if userCollection == nil {
		return nil, apperror.ErrUserCollectionNil
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := userCollection.Indexes().CreateOne(executionContext, indexModel)
	if err != nil {
		return nil, err
	}

	return &AuthRegisterRepository{
		userCollection: userCollection,
		userDocumentInserter: mongoUserDocumentInserter{
			userCollection: userCollection,
		},
	}, nil
}

func (repository *AuthRegisterRepository) CreateRegisterUser(executionContext context.Context, request *repositorymodel.CreateRegisterUserRequestModel, response *repositorymodel.CreateRegisterUserModel) error {
	if repository.userDocumentInserter == nil {
		return apperror.ErrAuthRegisterRepositoryNotConfigured
	}

	createdAt := datetime.GetCurrentDateTimeNow()
	document := repositorymodel.CreateRegisterUserDocumentModel{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
		CreatedAt:    createdAt,
	}

	insertedID, err := repository.userDocumentInserter.InsertOne(executionContext, document)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateKey
		}
		return err
	}

	response.ID = objectIDHex(insertedID)
	response.Name = request.Name
	response.Email = request.Email
	response.CreatedAt = createdAt
	return nil
}

func objectIDHex(insertedID any) string {
	objectID, ok := insertedID.(bson.ObjectID)
	if !ok {
		return ""
	}
	return objectID.Hex()
}
