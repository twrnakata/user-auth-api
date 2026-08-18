package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	repositorymodel "backend-challenge-golang-7solution/internal/repository/user/model"
)

type fakeUserDocumentFinder struct {
	called   bool
	filter   any
	document repositorymodel.GetUserByIDDocumentModel
	err      error
}

func (finder *fakeUserDocumentFinder) FindOne(executionContext context.Context, filter any, result any) error {
	finder.called = true
	finder.filter = filter
	if finder.err != nil {
		return finder.err
	}

	document, ok := result.(*repositorymodel.GetUserByIDDocumentModel)
	if !ok {
		return errors.New("unexpected get user document type")
	}
	*document = finder.document
	return nil
}

func TestNewGetUserRepository_NilCollection_ReturnsError(t *testing.T) {
	repository, err := NewGetUserRepository(nil)
	require.Error(t, err)
	require.Nil(t, repository)
}

func TestGetUserRepository_GetUserByID_FillsUserFromDocument(t *testing.T) {
	objectID := bson.NewObjectID()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	finder := &fakeUserDocumentFinder{
		document: repositorymodel.GetUserByIDDocumentModel{
			ID:        objectID,
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: createdAt,
		},
	}
	repository := &GetUserRepository{
		userDocumentFinder: finder,
	}

	var user domainuser.User
	err := repository.GetUserByID(context.Background(), objectID.Hex(), &user)

	require.NoError(t, err)
	require.True(t, finder.called)
	require.Equal(t, bson.M{"_id": objectID}, finder.filter)
	require.Equal(t, objectID.Hex(), user.ID)
	require.Equal(t, "Alice", user.Name)
	require.Equal(t, "alice@example.com", user.Email)
	require.Equal(t, createdAt, user.CreatedAt)
}

func TestGetUserRepository_GetUserByID_InvalidObjectID_ReturnsErrInvalidObjectID(t *testing.T) {
	finder := &fakeUserDocumentFinder{}
	repository := &GetUserRepository{
		userDocumentFinder: finder,
	}

	err := repository.GetUserByID(context.Background(), "not-an-object-id", &domainuser.User{})

	require.ErrorIs(t, err, ErrInvalidObjectID)
	require.False(t, finder.called)
}

func TestGetUserRepository_GetUserByID_NotFound_ReturnsErrNotFound(t *testing.T) {
	objectID := bson.NewObjectID()
	finder := &fakeUserDocumentFinder{
		err: mongo.ErrNoDocuments,
	}
	repository := &GetUserRepository{
		userDocumentFinder: finder,
	}

	err := repository.GetUserByID(context.Background(), objectID.Hex(), &domainuser.User{})

	require.ErrorIs(t, err, ErrNotFound)
	require.True(t, finder.called)
}

func TestGetUserRepository_GetUserByID_OtherFindError_ReturnsSameError(t *testing.T) {
	objectID := bson.NewObjectID()
	findError := errors.New("network timeout")
	finder := &fakeUserDocumentFinder{
		err: findError,
	}
	repository := &GetUserRepository{
		userDocumentFinder: finder,
	}

	err := repository.GetUserByID(context.Background(), objectID.Hex(), &domainuser.User{})

	require.ErrorIs(t, err, findError)
	require.False(t, errors.Is(err, ErrNotFound))
	require.False(t, errors.Is(err, ErrInvalidObjectID))
}
