package user

import (
	"context"
	"errors"

	repositoryuser "backend-challenge-golang-7solution/internal/repository/user"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
)

type DeleteUserService struct {
	Repository domainuser.DeleteUserRepository
}

func (service *DeleteUserService) Delete(executionContext context.Context, userID string, user *domainuser.User) error {
	if service.Repository == nil {
		return errors.New("delete user repository not configured")
	}
	if user == nil {
		return errors.New("user response is nil")
	}

	err := service.Repository.DeleteUser(executionContext, userID, user)
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
