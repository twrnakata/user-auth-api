package auth

import (
	"context"

	repositorymodel "github.com/twrnakata/user-auth-api/internal/repository/auth/model"
)

type AuthLoginRepository interface {
	GetLoginUserByEmail(executionContext context.Context, request *repositorymodel.AuthLoginRepositoryRequestModel, response *repositorymodel.GetLoginUserByEmailModel) error
}
