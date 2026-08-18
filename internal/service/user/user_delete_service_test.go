package user

import (
	"context"
	"errors"
	"testing"
	"time"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	repositoryuser "backend-challenge-golang-7solution/internal/repository/user"
)

type fakeDeleteUserRepository struct {
	called bool
	gotID  string
	user   *domainuser.User
	err    error
}

func (repository *fakeDeleteUserRepository) DeleteUser(executionContext context.Context, userID string, user *domainuser.User) error {
	repository.called = true
	repository.gotID = userID
	if user != nil && repository.user != nil {
		*user = *repository.user
	}
	return repository.err
}

func TestDeleteUserService_Delete_FillsUserFromRepository(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	repository := &fakeDeleteUserRepository{
		user: &domainuser.User{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: createdAt,
		},
	}

	service := &DeleteUserService{
		Repository: repository,
	}

	var user domainuser.User
	err := service.Delete(context.Background(), "u-1", &user)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !repository.called {
		t.Fatalf("expected repository DeleteUser to be called")
	}
	if repository.gotID != "u-1" {
		t.Fatalf("unexpected id passed to repository: %s", repository.gotID)
	}
	if user.ID != "u-1" || user.Name != "Alice" || user.Email != "alice@example.com" || !user.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestDeleteUserService_Delete_NotFound_ReturnsErrUserNotFound(t *testing.T) {
	repository := &fakeDeleteUserRepository{
		err: repositoryuser.ErrNotFound,
	}

	service := &DeleteUserService{
		Repository: repository,
	}

	err := service.Delete(context.Background(), "missing", &domainuser.User{})
	if !errors.Is(err, domainuser.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestDeleteUserService_Delete_InvalidObjectID_ReturnsErrInvalidUserID(t *testing.T) {
	repository := &fakeDeleteUserRepository{
		err: repositoryuser.ErrInvalidObjectID,
	}

	service := &DeleteUserService{
		Repository: repository,
	}

	err := service.Delete(context.Background(), "not-an-object-id", &domainuser.User{})
	if !errors.Is(err, domainuser.ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got: %v", err)
	}
}
