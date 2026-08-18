package auth

import (
	"context"
	"errors"

	servicemodel "backend-challenge-golang/internal/service/auth/model"
)

var ErrNotImplemented = errors.New("not implemented")

// RegisterUserService is a port for the Register use case.
// Handler depends on this interface to keep TDD simple.
type RegisterUserService interface {
	Register(executionContext context.Context, request *servicemodel.RegisterUserRequest, response *servicemodel.RegisterUserResponse) error
}

// NotImplementedRegisterUserService is used by main.go for now.
// Unit tests for the handler will inject their own fakes/mocks.
type NotImplementedRegisterUserService struct{}

func (service *NotImplementedRegisterUserService) Register(executionContext context.Context, request *servicemodel.RegisterUserRequest, response *servicemodel.RegisterUserResponse) error {
	return ErrNotImplemented
}
