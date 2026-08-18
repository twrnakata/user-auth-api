package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	repositorymodel "github.com/twrnakata/user-auth-api/internal/repository/user/model"
)

type fakeUserDocumentDeleter struct {
	called   bool
	filter   any
	document repositorymodel.DeleteUserDocumentModel
	err      error
}

func (deleter *fakeUserDocumentDeleter) FindOneAndDelete(executionContext context.Context, filter any, result any) error {
	deleter.called = true
	deleter.filter = filter
	if deleter.err != nil {
		return deleter.err
	}

	document, ok := result.(*repositorymodel.DeleteUserDocumentModel)
	if !ok {
		return errors.New("unexpected delete user document type")
	}
	*document = deleter.document
	return nil
}

func TestNewDeleteUserRepository_NilCollection_ReturnsError(t *testing.T) {
	repository, err := NewDeleteUserRepository(nil)
	require.Error(t, err)
	require.Nil(t, repository)
}

func TestDeleteUserRepository_DeleteUser_FillsUserFromDocument(t *testing.T) {
	objectID := bson.NewObjectID()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	deleter := &fakeUserDocumentDeleter{
		document: repositorymodel.DeleteUserDocumentModel{
			ID:        objectID,
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: createdAt,
		},
	}
	repository := &DeleteUserRepository{
		userDocumentDeleter: deleter,
	}

	var user domainuser.User
	err := repository.DeleteUser(context.Background(), objectID.Hex(), &user)

	require.NoError(t, err)
	require.True(t, deleter.called)
	require.Equal(t, bson.M{"_id": objectID}, deleter.filter)
	require.Equal(t, objectID.Hex(), user.ID)
	require.Equal(t, "Alice", user.Name)
	require.Equal(t, "alice@example.com", user.Email)
	require.Equal(t, createdAt, user.CreatedAt)
}

func TestDeleteUserRepository_DeleteUser_InvalidObjectID_ReturnsErrInvalidObjectID(t *testing.T) {
	deleter := &fakeUserDocumentDeleter{}
	repository := &DeleteUserRepository{
		userDocumentDeleter: deleter,
	}

	err := repository.DeleteUser(context.Background(), "not-an-object-id", &domainuser.User{})

	require.ErrorIs(t, err, ErrInvalidObjectID)
	require.False(t, deleter.called)
}

func TestDeleteUserRepository_DeleteUser_NotFound_ReturnsErrNotFound(t *testing.T) {
	objectID := bson.NewObjectID()
	deleter := &fakeUserDocumentDeleter{
		err: mongo.ErrNoDocuments,
	}
	repository := &DeleteUserRepository{
		userDocumentDeleter: deleter,
	}

	err := repository.DeleteUser(context.Background(), objectID.Hex(), &domainuser.User{})

	require.ErrorIs(t, err, ErrNotFound)
	require.True(t, deleter.called)
}

func TestDeleteUserRepository_DeleteUser_OtherDeleteError_ReturnsSameError(t *testing.T) {
	objectID := bson.NewObjectID()
	deleteError := errors.New("network timeout")
	deleter := &fakeUserDocumentDeleter{
		err: deleteError,
	}
	repository := &DeleteUserRepository{
		userDocumentDeleter: deleter,
	}

	err := repository.DeleteUser(context.Background(), objectID.Hex(), &domainuser.User{})

	require.ErrorIs(t, err, deleteError)
	require.False(t, errors.Is(err, ErrNotFound))
	require.False(t, errors.Is(err, ErrInvalidObjectID))
}
