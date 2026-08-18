package user

import (
	"context"
	"errors"

	repositoryuser "backend-challenge-golang-7solution/internal/repository/user"
	repositorymodel "backend-challenge-golang-7solution/internal/repository/user/model"
	servicemodel "backend-challenge-golang-7solution/internal/service/user/model"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
)

type GetUserService struct {
	Repository domainuser.GetUserRepository
}

func (service *GetUserService) GetByID(executionContext context.Context, request *servicemodel.GetUserRequestModel, response *servicemodel.GetUserResponseModel) error {
	if service.Repository == nil {
		return errors.New("get user repository not configured")
	}

	getUserByIDModel := &repositorymodel.GetUserByIDModel{}
	err := service.Repository.GetUserByID(executionContext, request, getUserByIDModel)
	if err != nil {
		if errors.Is(err, repositoryuser.ErrInvalidObjectID) {
			return domainuser.ErrInvalidUserID
		}
		if errors.Is(err, repositoryuser.ErrNotFound) {
			return domainuser.ErrUserNotFound
		}
		return err
	}

	response.ID = getUserByIDModel.ID
	response.Name = getUserByIDModel.Name
	response.Email = getUserByIDModel.Email
	response.CreatedAt = getUserByIDModel.CreatedAt
	return nil
}
