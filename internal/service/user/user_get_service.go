package user

import (
	"context"
	"errors"

	repositoryuser "github.com/twrnakata/user-auth-api/internal/repository/user"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

type GetUserService struct {
	Repository domainuser.GetUserRepository
}

func (service *GetUserService) GetByID(executionContext context.Context, userID string, user *domainuser.User) error {
	if service.Repository == nil {
		return apperror.ErrGetUserRepositoryNotConfigured
	}
	if user == nil {
		return apperror.ErrUserResponseNil
	}

	err := service.Repository.GetUserByID(executionContext, userID, user)
	if err != nil {
		if errors.Is(err, repositoryuser.ErrInvalidObjectID) {
			return domainuser.ErrInvalidUserID
		}
		if errors.Is(err, repositoryuser.ErrNotFound) {
			return domainuser.ErrUserNotFound
		}
		return err
	}

	return nil
}
