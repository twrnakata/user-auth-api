package auth

import (
	"context"
	"testing"
	"time"

	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type fakeRegisterUserRepository struct {
	called   bool
	gotReq   repositorymodel.CreateRegisterUserRequest
	response *repositorymodel.CreateRegisterUserResponse
	err      error
}

func (repository *fakeRegisterUserRepository) Create(executionContext context.Context, request *repositorymodel.CreateRegisterUserRequest, response *repositorymodel.CreateRegisterUserResponse) error {
	repository.called = true
	if request != nil {
		repository.gotReq = *request
	}
	if response != nil && repository.response != nil {
		*response = *repository.response
	}
	return repository.err
}

func TestRegisterUserService_HashesPasswordAndCallsRepository(t *testing.T) {
	repository := &fakeRegisterUserRepository{
		response: &repositorymodel.CreateRegisterUserResponse{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}

	service := &RegisterUserService{
		Repo: repository,
	}

	var response servicemodel.RegisterUserResponse
	err := service.Register(context.Background(), &servicemodel.RegisterUserRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret",
	}, &response)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !repository.called {
		t.Fatalf("expected repository Create to be called")
	}
	if repository.gotReq.PasswordHash == "" {
		t.Fatalf("expected password hash to be set")
	}
	if repository.gotReq.PasswordHash == "secret" {
		t.Fatalf("expected password hash to not equal plain password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repository.gotReq.PasswordHash), []byte("secret")); err != nil {
		t.Fatalf("expected stored hash to match plain password, compare error: %v", err)
	}

	if response.ID != "u-1" || response.Name != "Alice" || response.Email != "alice@example.com" {
		t.Fatalf("unexpected response: %#v", response)
	}
}
