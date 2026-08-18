package auth

import (
	"context"
	"errors"
	"testing"

	domainauth "backend-challenge-golang/internal/domain/auth"
	repositoryauth "backend-challenge-golang/internal/repository/auth"
	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type fakeAuthLoginRepository struct {
	called   bool
	gotEmail string
	response *repositorymodel.AuthLoginRepositoryResponse
	err      error
}

func (repository *fakeAuthLoginRepository) GetLoginUserByEmail(executionContext context.Context, request *repositorymodel.AuthLoginRepositoryRequest, response *repositorymodel.AuthLoginRepositoryResponse) error {
	repository.called = true

	if request != nil {
		repository.gotEmail = request.Email
	}

	if response != nil && repository.response != nil {
		*response = *repository.response
	}

	return repository.err
}

func TestAuthLoginService_Login_FillsResponseFromRepository(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("expected no error generating password hash, got: %v", err)
	}

	repository := &fakeAuthLoginRepository{
		response: &repositorymodel.AuthLoginRepositoryResponse{
			ID:           "u-1",
			Name:         "Alice",
			Email:        "alice@example.com",
			PasswordHash: string(passwordHash),
		},
	}

	service := &AuthLoginService{
		Repository: repository,
		BuildToken: func(userID string, name string) string {
			return "mock-token"
		},
	}

	request := &servicemodel.LoginUserRequest{
		Email:    "alice@example.com",
		Password: "secret",
	}
	response := &servicemodel.LoginUserResponse{}

	err = service.Login(context.Background(), request, response)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !repository.called {
		t.Fatalf("expected repository to be called")
	}

	if repository.gotEmail != "alice@example.com" {
		t.Fatalf("unexpected email passed to repository: %s", repository.gotEmail)
	}

	if response.Token != "mock-token" || response.ID != "u-1" || response.Name != "Alice" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestAuthLoginService_Login_NotFound_ReturnsErrUserNotFound(t *testing.T) {
	repository := &fakeAuthLoginRepository{
		err: repositoryauth.ErrNotFound,
	}

	service := &AuthLoginService{
		Repository: repository,
	}

	err := service.Login(context.Background(), &servicemodel.LoginUserRequest{
		Email:    "unknown@example.com",
		Password: "secret",
	}, &servicemodel.LoginUserResponse{})
	if !errors.Is(err, domainauth.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

