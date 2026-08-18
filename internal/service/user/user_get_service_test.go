package user

import (
	"context"
	"errors"
	"testing"
	"time"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	repositoryuser "github.com/twrnakata/user-auth-api/internal/repository/user"
)

type fakeGetUserRepository struct {
	called bool
	gotID  string
	user   *domainuser.User
	err    error
}

func (repository *fakeGetUserRepository) GetUserByID(executionContext context.Context, userID string, user *domainuser.User) error {
	repository.called = true
	repository.gotID = userID
	if user != nil && repository.user != nil {
		*user = *repository.user
	}
	return repository.err
}

func TestGetUserService_GetByID_FillsUserFromRepository(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	repository := &fakeGetUserRepository{
		user: &domainuser.User{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: createdAt,
		},
	}

	service := &GetUserService{
		Repository: repository,
	}

	var user domainuser.User
	err := service.GetByID(context.Background(), "u-1", &user)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !repository.called {
		t.Fatalf("expected repository GetUserByID to be called")
	}
	if repository.gotID != "u-1" {
		t.Fatalf("unexpected id passed to repository: %s", repository.gotID)
	}
	if user.ID != "u-1" || user.Name != "Alice" || user.Email != "alice@example.com" || !user.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestGetUserService_GetByID_NotFound_ReturnsErrUserNotFound(t *testing.T) {
	repository := &fakeGetUserRepository{
		err: repositoryuser.ErrNotFound,
	}

	service := &GetUserService{
		Repository: repository,
	}

	err := service.GetByID(context.Background(), "missing", &domainuser.User{})
	if !errors.Is(err, domainuser.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestGetUserService_GetByID_InvalidObjectID_ReturnsErrInvalidUserID(t *testing.T) {
	repository := &fakeGetUserRepository{
		err: repositoryuser.ErrInvalidObjectID,
	}

	service := &GetUserService{
		Repository: repository,
	}

	err := service.GetByID(context.Background(), "not-an-object-id", &domainuser.User{})
	if !errors.Is(err, domainuser.ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got: %v", err)
	}
}
