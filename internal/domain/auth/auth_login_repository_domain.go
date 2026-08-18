package auth

import (
	"context"

	repositorymodel "backend-challenge-golang-7solution/internal/repository/auth/model"
)

type AuthLoginRepository interface {
	GetLoginUserByEmail(executionContext context.Context, request *repositorymodel.AuthLoginRepositoryRequestModel, response *repositorymodel.GetLoginUserByEmailModel) error
}
