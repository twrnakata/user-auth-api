package user

import (
	"context"
	"errors"
	"testing"
	"time"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	repositoryuser "github.com/twrnakata/user-auth-api/internal/repository/user"
)

type fakeUpdateUserRepository struct {
	called    bool
	gotID     string
	gotUpdate domainuser.UserUpdate
	user      *domainuser.User
	err       error
}

func (repository *fakeUpdateUserRepository) UpdateUser(executionContext context.Context, userID string, update domainuser.UserUpdate, user *domainuser.User) error {
	repository.called = true
	repository.gotID = userID
	repository.gotUpdate = update
	if user != nil && repository.user != nil {
		*user = *repository.user
	}
	return repository.err
}

func TestUpdateUserService_Update_NameOnly_OmitsEmptyEmail(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	repository := &fakeUpdateUserRepository{
		user: &domainuser.User{
			ID:        "u-1",
			Name:      "Bob",
			Email:     "alice@example.com",
			CreatedAt: createdAt,
		},
	}

	service := &UpdateUserService{
		Repository: repository,
	}

	user := domainuser.User{
		Name: "Bob",
	}
	err := service.Update(context.Background(), "u-1", &user)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !repository.called {
		t.Fatalf("expected repository UpdateUser to be called")
	}
	if repository.gotID != "u-1" {
		t.Fatalf("unexpected id passed to repository: %s", repository.gotID)
	}
	if repository.gotUpdate.Name == nil || *repository.gotUpdate.Name != "Bob" {
		t.Fatalf("expected name patch, got: %#v", repository.gotUpdate.Name)
	}
	if repository.gotUpdate.Email != nil {
		t.Fatalf("expected email to be omitted, got: %#v", repository.gotUpdate.Email)
	}
	if user.ID != "u-1" || user.Name != "Bob" || user.Email != "alice@example.com" || !user.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestUpdateUserService_Update_EmailOnly_OmitsEmptyName(t *testing.T) {
	repository := &fakeUpdateUserRepository{
		user: &domainuser.User{
			ID:    "u-1",
			Name:  "Alice",
			Email: "bob@example.com",
		},
	}

	service := &UpdateUserService{
		Repository: repository,
	}

	user := domainuser.User{
		Email: "bob@example.com",
	}
	err := service.Update(context.Background(), "u-1", &user)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if repository.gotUpdate.Email == nil || *repository.gotUpdate.Email != "bob@example.com" {
		t.Fatalf("expected email patch, got: %#v", repository.gotUpdate.Email)
	}
	if repository.gotUpdate.Name != nil {
		t.Fatalf("expected name to be omitted, got: %#v", repository.gotUpdate.Name)
	}
}

func TestUpdateUserService_Update_NotFound_ReturnsErrUserNotFound(t *testing.T) {
	repository := &fakeUpdateUserRepository{
		err: repositoryuser.ErrNotFound,
	}

	service := &UpdateUserService{
		Repository: repository,
	}

	err := service.Update(context.Background(), "missing", &domainuser.User{Name: "Bob"})
	if !errors.Is(err, domainuser.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUpdateUserService_Update_InvalidObjectID_ReturnsErrInvalidUserID(t *testing.T) {
	repository := &fakeUpdateUserRepository{
		err: repositoryuser.ErrInvalidObjectID,
	}

	service := &UpdateUserService{
		Repository: repository,
	}

	err := service.Update(context.Background(), "not-an-object-id", &domainuser.User{Name: "Bob"})
	if !errors.Is(err, domainuser.ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got: %v", err)
	}
}

func TestUpdateUserService_Update_DuplicateKey_ReturnsErrEmailAlreadyExists(t *testing.T) {
	repository := &fakeUpdateUserRepository{
		err: repositoryuser.ErrDuplicateKey,
	}

	service := &UpdateUserService{
		Repository: repository,
	}

	err := service.Update(context.Background(), "u-1", &domainuser.User{Email: "taken@example.com"})
	if !errors.Is(err, domainuser.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got: %v", err)
	}
}
