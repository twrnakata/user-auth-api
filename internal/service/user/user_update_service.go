package user

import (
	"context"
	"errors"

	repositoryuser "backend-challenge-golang-7solution/internal/repository/user"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
)

type UpdateUserService struct {
	Repository domainuser.UpdateUserRepository
}

func (service *UpdateUserService) Update(executionContext context.Context, userID string, user *domainuser.User) error {
	if service.Repository == nil {
		return errors.New("update user repository not configured")
	}
	if user == nil {
		return errors.New("user response is nil")
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
