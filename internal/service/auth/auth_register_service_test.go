package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	domainauth "backend-challenge-golang-7solution/internal/domain/auth"
	repositoryauth "backend-challenge-golang-7solution/internal/repository/auth"
	repositorymodel "backend-challenge-golang-7solution/internal/repository/auth/model"
	servicemodel "backend-challenge-golang-7solution/internal/service/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type fakeRegisterUserRepository struct {
	called   bool
	gotReq   repositorymodel.CreateRegisterUserRequestModel
	response *repositorymodel.CreateRegisterUserModel
	err      error
}

func (repository *fakeRegisterUserRepository) CreateRegisterUser(executionContext context.Context, request *repositorymodel.CreateRegisterUserRequestModel, response *repositorymodel.CreateRegisterUserModel) error {
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
		response: &repositorymodel.CreateRegisterUserModel{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}

	service := &RegisterUserService{
		Repository: repository,
	}

	var response servicemodel.RegisterUserResponseModel
	err := service.Register(context.Background(), &servicemodel.RegisterUserRequestModel{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret",
	}, &response)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !repository.called {
		t.Fatalf("expected repository CreateRegisterUser to be called")
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

func TestRegisterUserService_DuplicateKey_ReturnsEmailAlreadyExists(t *testing.T) {
	repository := &fakeRegisterUserRepository{
		err: repositoryauth.ErrDuplicateKey,
	}

	service := &RegisterUserService{
		Repository: repository,
	}

	err := service.Register(context.Background(), &servicemodel.RegisterUserRequestModel{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret",
	}, &servicemodel.RegisterUserResponseModel{})
	if !errors.Is(err, domainauth.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got: %v", err)
	}
}
