package user

import (
	"context"
	"errors"
	"testing"
	"time"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
)

type fakeListUserRepository struct {
	called bool
	users  []domainuser.User
	err    error
}

func (repository *fakeListUserRepository) ListUsers(executionContext context.Context, users *[]domainuser.User) error {
	repository.called = true
	if users != nil && repository.users != nil {
		*users = repository.users
	}
	return repository.err
}

func TestListUserService_List_FillsUsersFromRepository(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	repository := &fakeListUserRepository{
		users: []domainuser.User{
			{
				ID:        "u-1",
				Name:      "Alice",
				Email:     "alice@example.com",
				CreatedAt: createdAt,
			},
		},
	}

	service := &ListUserService{
		Repository: repository,
	}

	var users []domainuser.User
	err := service.List(context.Background(), &users)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !repository.called {
		t.Fatalf("expected repository ListUsers to be called")
	}
	if len(users) != 1 {
		t.Fatalf("unexpected users: %#v", users)
	}
	if users[0].ID != "u-1" || users[0].Name != "Alice" || users[0].Email != "alice@example.com" || !users[0].CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected user: %#v", users[0])
	}
}

func TestListUserService_List_Empty_ReturnsEmptySlice(t *testing.T) {
	repository := &fakeListUserRepository{}

	service := &ListUserService{
		Repository: repository,
	}

	var users []domainuser.User
	err := service.List(context.Background(), &users)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if users == nil {
		t.Fatalf("expected empty users slice, got nil")
	}
	if len(users) != 0 {
		t.Fatalf("expected no users, got: %#v", users)
	}
}

func TestListUserService_List_RepositoryError_ReturnsSameError(t *testing.T) {
	repositoryError := errors.New("network timeout")
	repository := &fakeListUserRepository{
		err: repositoryError,
	}

	service := &ListUserService{
		Repository: repository,
	}

	var users []domainuser.User
	err := service.List(context.Background(), &users)
	if !errors.Is(err, repositoryError) {
		t.Fatalf("expected repository error, got: %v", err)
	}
}
