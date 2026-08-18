package user

import (
	"context"

	repositorymodel "backend-challenge-golang-7solution/internal/repository/user/model"
	servicemodel "backend-challenge-golang-7solution/internal/service/user/model"
)

type GetUserRepository interface {
	GetUserByID(executionContext context.Context, request *servicemodel.GetUserRequestModel, response *repositorymodel.GetUserByIDModel) error
}
