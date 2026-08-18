package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repositorymodel "backend-challenge-golang-7solution/internal/repository/auth/model"
)

type fakeUserDocumentFinder struct {
	called   bool
	filter   any
	document repositorymodel.GetLoginUserDocumentModel
	err      error
}

func (finder *fakeUserDocumentFinder) FindOne(executionContext context.Context, filter any, result any) error {
	finder.called = true
	finder.filter = filter
	if finder.err != nil {
		return finder.err
	}

	document, ok := result.(*repositorymodel.GetLoginUserDocumentModel)
	if !ok {
		return errors.New("unexpected login document type")
	}
	*document = finder.document
	return nil
}

func TestNewAuthLoginRepository_NilCollection_ReturnsError(t *testing.T) {
	repository, err := NewAuthLoginRepository(nil)
	require.Error(t, err)
	require.Nil(t, repository)
}

func TestAuthLoginRepository_GetLoginUserByEmail_FillsResponseFromDocument(t *testing.T) {
	objectID := bson.NewObjectID()
	finder := &fakeUserDocumentFinder{
		document: repositorymodel.GetLoginUserDocumentModel{
			ID:           objectID,
			Name:         "Alice",
			Email:        "alice@example.com",
			PasswordHash: "hashed-password",
		},
	}
	repository := &AuthLoginRepository{
		userDocumentFinder: finder,
	}

	var response repositorymodel.GetLoginUserByEmailModel
	err := repository.GetLoginUserByEmail(context.Background(), &repositorymodel.AuthLoginRepositoryRequestModel{
		Email: "alice@example.com",
	}, &response)

	require.NoError(t, err)
	require.True(t, finder.called)
	require.Equal(t, bson.M{"email": "alice@example.com"}, finder.filter)
	require.Equal(t, objectID.Hex(), response.ID)
	require.Equal(t, "Alice", response.Name)
	require.Equal(t, "alice@example.com", response.Email)
	require.Equal(t, "hashed-password", response.PasswordHash)
}

func TestAuthLoginRepository_GetLoginUserByEmail_NotFound_ReturnsErrNotFound(t *testing.T) {
	finder := &fakeUserDocumentFinder{
		err: mongo.ErrNoDocuments,
	}
	repository := &AuthLoginRepository{
		userDocumentFinder: finder,
	}

	err := repository.GetLoginUserByEmail(context.Background(), &repositorymodel.AuthLoginRepositoryRequestModel{
		Email: "unknown@example.com",
	}, &repositorymodel.GetLoginUserByEmailModel{})

	require.ErrorIs(t, err, ErrNotFound)
}

func TestAuthLoginRepository_GetLoginUserByEmail_OtherFindError_ReturnsSameError(t *testing.T) {
	findError := errors.New("network timeout")
	finder := &fakeUserDocumentFinder{
		err: findError,
	}
	repository := &AuthLoginRepository{
		userDocumentFinder: finder,
	}

	err := repository.GetLoginUserByEmail(context.Background(), &repositorymodel.AuthLoginRepositoryRequestModel{
		Email: "alice@example.com",
	}, &repositorymodel.GetLoginUserByEmailModel{})

	require.ErrorIs(t, err, findError)
	require.False(t, errors.Is(err, ErrNotFound))
}
