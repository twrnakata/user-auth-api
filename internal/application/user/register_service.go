package user

import (
	"context"
	"errors"
	"time"
)

var ErrNotImplemented = errors.New("not implemented")

type RegisterUserRequest struct {
	Name     string
	Email    string
	Password string
}

type RegisterUserResponse struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

// RegisterUserService is a port for the Register use case.
// Handler depends on this interface to keep TDD simple.
type RegisterUserService interface {
	Register(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error)
}

// NotImplementedRegisterUserService is used by main.go for now.
// Unit tests for the handler will inject their own fakes/mocks.
type NotImplementedRegisterUserService struct{}

func (s *NotImplementedRegisterUserService) Register(ctx context.Context, req RegisterUserRequest) (RegisterUserResponse, error) {
	return RegisterUserResponse{}, ErrNotImplemented
}

