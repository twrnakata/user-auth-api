package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
	"backend-challenge-golang/pkg/datetime"
)

type fakeUserDocumentInserter struct {
	called     bool
	document   any
	insertedID any
	err        error
}

func (inserter *fakeUserDocumentInserter) InsertOne(executionContext context.Context, document any) (any, error) {
	inserter.called = true
	inserter.document = document
	return inserter.insertedID, inserter.err
}

func TestNewAuthRegisterRepository_NilCollection_ReturnsError(t *testing.T) {
	repository, err := NewAuthRegisterRepository(context.Background(), nil)
	require.Error(t, err)
	require.Nil(t, repository)
}

func TestAuthRegisterRepository_CreateRegisterUser_InsertsDocumentAndFillsResponse(t *testing.T) {
	objectID := bson.NewObjectID()
	inserter := &fakeUserDocumentInserter{
		insertedID: objectID,
	}
	repository := &AuthRegisterRepository{
		userDocumentInserter: inserter,
	}

	var response repositorymodel.CreateRegisterUserModel
	beforeInsert := datetime.GetCurrentDateTimeNow().Add(-time.Second)
	err := repository.CreateRegisterUser(context.Background(), &repositorymodel.CreateRegisterUserRequestModel{
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hashed-password",
	}, &response)
	afterInsert := datetime.GetCurrentDateTimeNow().Add(time.Second)

	require.NoError(t, err)
	require.True(t, inserter.called)

	document, ok := inserter.document.(repositorymodel.CreateRegisterUserDocumentModel)
	require.True(t, ok)
	require.Equal(t, "Alice", document.Name)
	require.Equal(t, "alice@example.com", document.Email)
	require.Equal(t, "hashed-password", document.PasswordHash)
	require.False(t, document.CreatedAt.IsZero())
	require.True(t, !document.CreatedAt.Before(beforeInsert) && !document.CreatedAt.After(afterInsert))

	require.Equal(t, objectID.Hex(), response.ID)
	require.Equal(t, "Alice", response.Name)
	require.Equal(t, "alice@example.com", response.Email)
	require.Equal(t, document.CreatedAt, response.CreatedAt)
}

func TestAuthRegisterRepository_CreateRegisterUser_DuplicateKey_ReturnsErrDuplicateKey(t *testing.T) {
	inserter := &fakeUserDocumentInserter{
		err: mongo.WriteException{
			WriteErrors: []mongo.WriteError{
				{Code: 11000, Message: "E11000 duplicate key error"},
			},
		},
	}
	repository := &AuthRegisterRepository{
		userDocumentInserter: inserter,
	}

	err := repository.CreateRegisterUser(context.Background(), &repositorymodel.CreateRegisterUserRequestModel{
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hashed-password",
	}, &repositorymodel.CreateRegisterUserModel{})

	require.ErrorIs(t, err, ErrDuplicateKey)
}

func TestAuthRegisterRepository_CreateRegisterUser_OtherInsertError_ReturnsSameError(t *testing.T) {
	insertError := errors.New("network timeout")
	inserter := &fakeUserDocumentInserter{
		err: insertError,
	}
	repository := &AuthRegisterRepository{
		userDocumentInserter: inserter,
	}

	err := repository.CreateRegisterUser(context.Background(), &repositorymodel.CreateRegisterUserRequestModel{
		Name:         "Alice",
		Email:        "alice@example.com",
		PasswordHash: "hashed-password",
	}, &repositorymodel.CreateRegisterUserModel{})

	require.ErrorIs(t, err, insertError)
	require.False(t, errors.Is(err, ErrDuplicateKey))
}
