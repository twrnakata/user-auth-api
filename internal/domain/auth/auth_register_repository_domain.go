package auth

import (
	"context"

	repositorymodel "backend-challenge-golang-7solution/internal/repository/auth/model"
)

// RegisterUserRepository is a port for persisting a new user during register.
// Domain keeps only interfaces (ports) so handler/service remain testable.
type RegisterUserRepository interface {
	CreateRegisterUser(executionContext context.Context, request *repositorymodel.CreateRegisterUserRequestModel, response *repositorymodel.CreateRegisterUserModel) error
}
