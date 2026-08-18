package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeUserDocumentsCounter struct {
	called bool
	filter any
	count  int64
	err    error
}

func (counter *fakeUserDocumentsCounter) CountDocuments(executionContext context.Context, filter any) (int64, error) {
	counter.called = true
	counter.filter = filter
	return counter.count, counter.err
}

func TestNewCountUserRepository_NilCollection_ReturnsError(t *testing.T) {
	repository, err := NewCountUserRepository(nil)
	require.Error(t, err)
	require.Nil(t, repository)
}

func TestCountUserRepository_CountUsers_FillsCountFromDocuments(t *testing.T) {
	counter := &fakeUserDocumentsCounter{
		count: 3,
	}
	repository := &CountUserRepository{
		userDocumentsCounter: counter,
	}

	var count int64
	err := repository.CountUsers(context.Background(), &count)

	require.NoError(t, err)
	require.True(t, counter.called)
	require.Equal(t, bson.M{}, counter.filter)
	require.Equal(t, int64(3), count)
}

func TestCountUserRepository_CountUsers_NilCount_ReturnsError(t *testing.T) {
	counter := &fakeUserDocumentsCounter{}
	repository := &CountUserRepository{
		userDocumentsCounter: counter,
	}

	err := repository.CountUsers(context.Background(), nil)

	require.Error(t, err)
	require.False(t, counter.called)
}

func TestCountUserRepository_CountUsers_CountError_ReturnsSameError(t *testing.T) {
	countError := errors.New("network timeout")
	counter := &fakeUserDocumentsCounter{
		err: countError,
	}
	repository := &CountUserRepository{
		userDocumentsCounter: counter,
	}

	var count int64
	err := repository.CountUsers(context.Background(), &count)

	require.ErrorIs(t, err, countError)
}
