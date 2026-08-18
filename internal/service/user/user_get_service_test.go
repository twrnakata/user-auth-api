package user

import (
	"context"
	"errors"
	"testing"
	"time"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	repositoryuser "backend-challenge-golang-7solution/internal/repository/user"
	repositorymodel "backend-challenge-golang-7solution/internal/repository/user/model"
	servicemodel "backend-challenge-golang-7solution/internal/service/user/model"
)

type fakeGetUserRepository struct {
	called   bool
	gotID    string
	response *repositorymodel.GetUserByIDModel
	err      error
}

func (repository *fakeGetUserRepository) GetUserByID(executionContext context.Context, request *servicemodel.GetUserRequestModel, response *repositorymodel.GetUserByIDModel) error {
	repository.called = true
	if request != nil {
		repository.gotID = request.ID
	}
	if response != nil && repository.response != nil {
		*response = *repository.response
	}
	return repository.err
}

func TestGetUserService_GetByID_FillsResponseFromRepository(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	repository := &fakeGetUserRepository{
		response: &repositorymodel.GetUserByIDModel{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: createdAt,
		},
	}

	service := &GetUserService{
		Repository: repository,
	}

	request := &servicemodel.GetUserRequestModel{
		ID: "u-1",
	}
	response := &servicemodel.GetUserResponseModel{}

	err := service.GetByID(context.Background(), request, response)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !repository.called {
		t.Fatalf("expected repository GetUserByID to be called")
	}
	if repository.gotID != "u-1" {
		t.Fatalf("unexpected id passed to repository: %s", repository.gotID)
	}

	if response.ID != "u-1" || response.Name != "Alice" || response.Email != "alice@example.com" || !response.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestGetUserService_GetByID_NotFound_ReturnsErrUserNotFound(t *testing.T) {
	repository := &fakeGetUserRepository{
		err: repositoryuser.ErrNotFound,
	}

	service := &GetUserService{
		Repository: repository,
	}

	err := service.GetByID(context.Background(), &servicemodel.GetUserRequestModel{
		ID: "missing",
	}, &servicemodel.GetUserResponseModel{})
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

	err := service.GetByID(context.Background(), &servicemodel.GetUserRequestModel{
		ID: "not-an-object-id",
	}, &servicemodel.GetUserResponseModel{})
	if !errors.Is(err, domainuser.ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got: %v", err)
	}
}
