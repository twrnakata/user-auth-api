package user

import "context"

type DeleteUserService interface {
	Delete(executionContext context.Context, userID string, user *User) error
}

type NotImplementedDeleteUserService struct{}

func (service *NotImplementedDeleteUserService) Delete(executionContext context.Context, userID string, user *User) error {
	return ErrNotImplemented
}
