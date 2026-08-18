package user

import (
	"context"
	"errors"

	repositoryuser "github.com/twrnakata/user-auth-api/internal/repository/user"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

type UpdateUserService struct {
	Repository domainuser.UpdateUserRepository
}

func (service *UpdateUserService) Update(executionContext context.Context, userID string, user *domainuser.User) error {
	if service.Repository == nil {
		return apperror.ErrUpdateUserRepositoryNotConfigured
	}
	if user == nil {
		return apperror.ErrUserResponseNil
	}

	update := domainuser.UserUpdate{}
	if user.Name != "" {
		name := user.Name
		update.Name = &name
	}
	if user.Email != "" {
		email := user.Email
		update.Email = &email
	}

	err := service.Repository.UpdateUser(executionContext, userID, update, user)
	if err != nil {
		if errors.Is(err, repositoryuser.ErrInvalidObjectID) {
			return domainuser.ErrInvalidUserID
		}
		if errors.Is(err, repositoryuser.ErrNotFound) {
			return domainuser.ErrUserNotFound
		}
		if errors.Is(err, repositoryuser.ErrDuplicateKey) {
			return domainuser.ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}
