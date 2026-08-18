package user

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	repositorymodel "github.com/twrnakata/user-auth-api/internal/repository/user/model"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

type userDocumentsFinder interface {
	Find(executionContext context.Context, filter any, sort any, result any) error
}

type mongoUserDocumentsFinder struct {
	userCollection *mongo.Collection
}

func (finder mongoUserDocumentsFinder) Find(executionContext context.Context, filter any, sort any, result any) error {
	cursor, err := finder.userCollection.Find(executionContext, filter, options.Find().SetSort(sort))
	if err != nil {
		return err
	}
	defer cursor.Close(executionContext)
	return cursor.All(executionContext, result)
}

type ListUserRepository struct {
	userCollection      *mongo.Collection
	userDocumentsFinder userDocumentsFinder
}

func NewListUserRepository(userCollection *mongo.Collection) (*ListUserRepository, error) {
	if userCollection == nil {
		return nil, apperror.ErrUserCollectionNil
	}

	return &ListUserRepository{
		userCollection: userCollection,
		userDocumentsFinder: mongoUserDocumentsFinder{
			userCollection: userCollection,
		},
	}, nil
}

func (repository *ListUserRepository) ListUsers(executionContext context.Context, users *[]domainuser.User) error {
	if repository.userDocumentsFinder == nil {
		return apperror.ErrListUserRepositoryNotConfigured
	}
	if users == nil {
		return apperror.ErrUsersResponseNil
	}

	var documents []repositorymodel.ListUsersDocumentModel
	err := repository.userDocumentsFinder.Find(executionContext, bson.M{}, bson.D{{Key: "createdAt", Value: -1}}, &documents)
	if err != nil {
		return err
	}

	listedUsers := make([]domainuser.User, 0, len(documents))
	for _, document := range documents {
		listedUsers = append(listedUsers, domainuser.User{
			ID:        document.ID.Hex(),
			Name:      document.Name,
			Email:     document.Email,
			CreatedAt: document.CreatedAt,
		})
	}
	*users = listedUsers
	return nil
}
