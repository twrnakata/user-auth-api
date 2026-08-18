package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	repositorymodel "github.com/twrnakata/user-auth-api/internal/repository/user/model"
)

type fakeUserDocumentsFinder struct {
	called    bool
	filter    any
	sort      any
	documents []repositorymodel.ListUsersDocumentModel
	err       error
}

func (finder *fakeUserDocumentsFinder) Find(executionContext context.Context, filter any, sort any, result any) error {
	finder.called = true
	finder.filter = filter
	finder.sort = sort
	if finder.err != nil {
		return finder.err
	}

	documents, ok := result.(*[]repositorymodel.ListUsersDocumentModel)
	if !ok {
		return errors.New("unexpected list user documents type")
	}
	*documents = finder.documents
	return nil
}

func TestNewListUserRepository_NilCollection_ReturnsError(t *testing.T) {
	repository, err := NewListUserRepository(nil)
	require.Error(t, err)
	require.Nil(t, repository)
}

func TestListUserRepository_ListUsers_FillsUsersFromDocuments(t *testing.T) {
	objectID := bson.NewObjectID()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	finder := &fakeUserDocumentsFinder{
		documents: []repositorymodel.ListUsersDocumentModel{
			{
				ID:        objectID,
				Name:      "Alice",
				Email:     "alice@example.com",
				CreatedAt: createdAt,
			},
		},
	}
	repository := &ListUserRepository{
		userDocumentsFinder: finder,
	}

	var users []domainuser.User
	err := repository.ListUsers(context.Background(), &users)

	require.NoError(t, err)
	require.True(t, finder.called)
	require.Equal(t, bson.M{}, finder.filter)
	require.Equal(t, bson.D{{Key: "createdAt", Value: -1}}, finder.sort)
	require.Len(t, users, 1)
	require.Equal(t, objectID.Hex(), users[0].ID)
	require.Equal(t, "Alice", users[0].Name)
	require.Equal(t, "alice@example.com", users[0].Email)
	require.Equal(t, createdAt, users[0].CreatedAt)
}

func TestListUserRepository_ListUsers_Empty_ReturnsEmptySlice(t *testing.T) {
	finder := &fakeUserDocumentsFinder{}
	repository := &ListUserRepository{
		userDocumentsFinder: finder,
	}

	var users []domainuser.User
	err := repository.ListUsers(context.Background(), &users)

	require.NoError(t, err)
	require.True(t, finder.called)
	require.NotNil(t, users)
	require.Empty(t, users)
}

func TestListUserRepository_ListUsers_FindError_ReturnsSameError(t *testing.T) {
	findError := errors.New("network timeout")
	finder := &fakeUserDocumentsFinder{
		err: findError,
	}
	repository := &ListUserRepository{
		userDocumentsFinder: finder,
	}

	var users []domainuser.User
	err := repository.ListUsers(context.Background(), &users)

	require.ErrorIs(t, err, findError)
}
