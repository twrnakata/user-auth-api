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

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidObjectID = errors.New("invalid object id")
	ErrDuplicateKey    = errors.New("duplicate key")
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
		return nil, apperror.ErrUserCollectionNil
	}

	return &GetUserRepository{
		userCollection: userCollection,
		userDocumentFinder: mongoUserDocumentFinder{
			userCollection: userCollection,
		},
	}, nil
}

func (repository *GetUserRepository) GetUserByID(executionContext context.Context, userID string, user *domainuser.User) error {
	if repository.userDocumentFinder == nil {
		return apperror.ErrGetUserRepositoryNotConfigured
	}
	if user == nil {
		return apperror.ErrUserResponseNil
	}

	objectID, err := bson.ObjectIDFromHex(userID)
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

	*user = domainuser.User{
		ID:        document.ID.Hex(),
		Name:      document.Name,
		Email:     document.Email,
		CreatedAt: document.CreatedAt,
	}
	return nil
}
