package user

import (
	"context"
	"errors"

	servicemodel "backend-challenge-golang/internal/service/user/model"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidUserID  = errors.New("invalid user id")
)

type GetUserService interface {
	GetByID(executionContext context.Context, request *servicemodel.GetUserRequestModel, response *servicemodel.GetUserResponseModel) error
}

type NotImplementedGetUserService struct{}

func (service *NotImplementedGetUserService) GetByID(executionContext context.Context, request *servicemodel.GetUserRequestModel, response *servicemodel.GetUserResponseModel) error {
	return ErrNotImplemented
}
