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
		return nil, errors.New("user collection is nil")
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
		return errors.New("list user repository not configured")
	}
	if users == nil {
		return errors.New("users response is nil")
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
