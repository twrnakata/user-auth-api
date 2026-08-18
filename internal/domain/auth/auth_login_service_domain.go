package auth

import (
	"context"
	"errors"

	servicemodel "github.com/twrnakata/user-auth-api/internal/service/auth/model"
)

var (
	// ErrInvalidCredentials is returned for both unknown email and wrong password.
	// Login must not distinguish those cases (no 404 vs 401) so a caller cannot
	// probe which emails are registered.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// LoginUserService is a port for the login use case.
type LoginUserService interface {
	Login(executionContext context.Context, request *servicemodel.LoginUserRequestModel, response *servicemodel.LoginUserResponseModel) error
}

type NotImplementedLoginUserService struct{}

func (service *NotImplementedLoginUserService) Login(executionContext context.Context, request *servicemodel.LoginUserRequestModel, response *servicemodel.LoginUserResponseModel) error {
	return ErrNotImplemented
}
