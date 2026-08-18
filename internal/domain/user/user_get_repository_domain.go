package user

import (
	"context"

	repositorymodel "backend-challenge-golang/internal/repository/user/model"
	servicemodel "backend-challenge-golang/internal/service/user/model"
)

type GetUserRepository interface {
	GetUserByID(executionContext context.Context, request *servicemodel.GetUserRequestModel, response *repositorymodel.GetUserByIDModel) error
}
