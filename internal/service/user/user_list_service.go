package user

import (
	"context"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

type ListUserService struct {
	Repository domainuser.ListUserRepository
}

func (service *ListUserService) List(executionContext context.Context, users *[]domainuser.User) error {
	if service.Repository == nil {
		return apperror.ErrListUserRepositoryNotConfigured
	}
	if users == nil {
		return apperror.ErrUsersResponseNil
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
