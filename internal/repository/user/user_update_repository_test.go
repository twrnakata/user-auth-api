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

type fakeUserDocumentUpdater struct {
	called   bool
	filter   any
	update   any
	document repositorymodel.UpdateUserDocumentModel
	err      error
}

func (updater *fakeUserDocumentUpdater) FindOneAndUpdate(executionContext context.Context, filter any, update any, result any) error {
	updater.called = true
	updater.filter = filter
	updater.update = update
	if updater.err != nil {
		return updater.err
	}

	document, ok := result.(*repositorymodel.UpdateUserDocumentModel)
	if !ok {
		return errors.New("unexpected update user document type")
	}
	*document = updater.document
	return nil
}

func TestNewUpdateUserRepository_NilCollection_ReturnsError(t *testing.T) {
	repository, err := NewUpdateUserRepository(nil)
	require.Error(t, err)
	require.Nil(t, repository)
}

func stringPointer(value string) *string {
	return &value
}

func TestUpdateUserRepository_UpdateUser_NameOnly_FillsUserFromDocument(t *testing.T) {
	objectID := bson.NewObjectID()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updater := &fakeUserDocumentUpdater{
		document: repositorymodel.UpdateUserDocumentModel{
			ID:        objectID,
			Name:      "Bob",
			Email:     "alice@example.com",
			CreatedAt: createdAt,
		},
	}
	repository := &UpdateUserRepository{
		userDocumentUpdater: updater,
	}

	var user domainuser.User
	err := repository.UpdateUser(context.Background(), objectID.Hex(), domainuser.UserUpdate{
		Name: stringPointer("Bob"),
	}, &user)

	require.NoError(t, err)
	require.True(t, updater.called)
	require.Equal(t, bson.M{"_id": objectID}, updater.filter)
	require.Equal(t, bson.M{"$set": bson.M{"name": "Bob"}}, updater.update)
	require.Equal(t, objectID.Hex(), user.ID)
	require.Equal(t, "Bob", user.Name)
	require.Equal(t, "alice@example.com", user.Email)
	require.Equal(t, createdAt, user.CreatedAt)
}

func TestUpdateUserRepository_UpdateUser_EmailOnly_SetsEmailOnly(t *testing.T) {
	objectID := bson.NewObjectID()
	updater := &fakeUserDocumentUpdater{
		document: repositorymodel.UpdateUserDocumentModel{
			ID:    objectID,
			Name:  "Alice",
			Email: "bob@example.com",
		},
	}
	repository := &UpdateUserRepository{
		userDocumentUpdater: updater,
	}

	var user domainuser.User
	err := repository.UpdateUser(context.Background(), objectID.Hex(), domainuser.UserUpdate{
		Email: stringPointer("bob@example.com"),
	}, &user)

	require.NoError(t, err)
	require.Equal(t, bson.M{"$set": bson.M{"email": "bob@example.com"}}, updater.update)
	require.Equal(t, "Alice", user.Name)
	require.Equal(t, "bob@example.com", user.Email)
}

func TestUpdateUserRepository_UpdateUser_InvalidObjectID_ReturnsErrInvalidObjectID(t *testing.T) {
	updater := &fakeUserDocumentUpdater{}
	repository := &UpdateUserRepository{
		userDocumentUpdater: updater,
	}

	err := repository.UpdateUser(context.Background(), "not-an-object-id", domainuser.UserUpdate{
		Name: stringPointer("Bob"),
	}, &domainuser.User{})

	require.ErrorIs(t, err, ErrInvalidObjectID)
	require.False(t, updater.called)
}

func TestUpdateUserRepository_UpdateUser_NotFound_ReturnsErrNotFound(t *testing.T) {
	objectID := bson.NewObjectID()
	updater := &fakeUserDocumentUpdater{
		err: mongo.ErrNoDocuments,
	}
	repository := &UpdateUserRepository{
		userDocumentUpdater: updater,
	}

	err := repository.UpdateUser(context.Background(), objectID.Hex(), domainuser.UserUpdate{
		Name: stringPointer("Bob"),
	}, &domainuser.User{})

	require.ErrorIs(t, err, ErrNotFound)
	require.True(t, updater.called)
}

func TestUpdateUserRepository_UpdateUser_DuplicateKey_ReturnsErrDuplicateKey(t *testing.T) {
	objectID := bson.NewObjectID()
	updater := &fakeUserDocumentUpdater{
		err: mongo.WriteException{
			WriteErrors: []mongo.WriteError{
				{Code: 11000, Message: "E11000 duplicate key error"},
			},
		},
	}
	repository := &UpdateUserRepository{
		userDocumentUpdater: updater,
	}

	err := repository.UpdateUser(context.Background(), objectID.Hex(), domainuser.UserUpdate{
		Email: stringPointer("taken@example.com"),
	}, &domainuser.User{})

	require.ErrorIs(t, err, ErrDuplicateKey)
}

func TestUpdateUserRepository_UpdateUser_OtherUpdateError_ReturnsSameError(t *testing.T) {
	objectID := bson.NewObjectID()
	updateError := errors.New("network timeout")
	updater := &fakeUserDocumentUpdater{
		err: updateError,
	}
	repository := &UpdateUserRepository{
		userDocumentUpdater: updater,
	}

	err := repository.UpdateUser(context.Background(), objectID.Hex(), domainuser.UserUpdate{
		Name: stringPointer("Bob"),
	}, &domainuser.User{})

	require.ErrorIs(t, err, updateError)
	require.False(t, errors.Is(err, ErrNotFound))
	require.False(t, errors.Is(err, ErrDuplicateKey))
}
