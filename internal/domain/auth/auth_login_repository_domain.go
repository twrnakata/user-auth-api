package auth

import (
	"context"

	repositorymodel "backend-challenge-golang/internal/repository/auth/model"
)

type AuthLoginRepository interface {
	GetLoginUserByEmail(executionContext context.Context, request *repositorymodel.AuthLoginRepositoryRequest, response *repositorymodel.AuthLoginRepositoryResponse) error
}

