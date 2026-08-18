package auth

import (
	"context"
	"errors"

	servicemodel "backend-challenge-golang-7solution/internal/service/auth/model"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

// LoginUserService is a port for the login use case.
type LoginUserService interface {
	Login(executionContext context.Context, request *servicemodel.LoginUserRequestModel, response *servicemodel.LoginUserResponseModel) error
}

type NotImplementedLoginUserService struct{}

func (service *NotImplementedLoginUserService) Login(executionContext context.Context, request *servicemodel.LoginUserRequestModel, response *servicemodel.LoginUserResponseModel) error {
	return ErrNotImplemented
}
