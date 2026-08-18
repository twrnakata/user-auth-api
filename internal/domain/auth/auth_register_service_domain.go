package auth

import (
	"context"
	"errors"

	servicemodel "backend-challenge-golang-7solution/internal/service/auth/model"
)

var (
	ErrNotImplemented     = errors.New("not implemented")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// RegisterUserService is a port for the Register use case.
// Handler depends on this interface to keep TDD simple.
type RegisterUserService interface {
	Register(executionContext context.Context, request *servicemodel.RegisterUserRequestModel, response *servicemodel.RegisterUserResponseModel) error
}

// NotImplementedRegisterUserService is used by main.go for now.
// Unit tests for the handler will inject their own fakes/mocks.
type NotImplementedRegisterUserService struct{}

func (service *NotImplementedRegisterUserService) Register(executionContext context.Context, request *servicemodel.RegisterUserRequestModel, response *servicemodel.RegisterUserResponseModel) error {
	return ErrNotImplemented
}
