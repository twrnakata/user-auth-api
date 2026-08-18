package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	repositorymodel "backend-challenge-golang-7solution/internal/repository/user/model"
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
		return nil, errors.New("user collection is nil")
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
		return errors.New("get user repository not configured")
	}
	if user == nil {
		return errors.New("user response is nil")
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
