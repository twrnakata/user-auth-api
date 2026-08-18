package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnect_EmptyURI_ReturnsError(t *testing.T) {
	client, err := Connect(context.Background(), "")
	require.Error(t, err)
	require.Nil(t, client)
}

func TestUserCollection_UsesUsersCollection(t *testing.T) {
	require.Equal(t, "users", UsersCollectionName)
}
