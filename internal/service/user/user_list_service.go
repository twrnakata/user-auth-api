package user

import (
	"context"
	"errors"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
)

type ListUserService struct {
	Repository domainuser.ListUserRepository
}

func (service *ListUserService) List(executionContext context.Context, users *[]domainuser.User) error {
	if service.Repository == nil {
		return errors.New("list user repository not configured")
	}
	if users == nil {
		return errors.New("users response is nil")
	}

	err := service.Repository.ListUsers(executionContext, users)
	if err != nil {
		return err
	}

	if *users == nil {
		*users = []domainuser.User{}
	}
	return nil
}
