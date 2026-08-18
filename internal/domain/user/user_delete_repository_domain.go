package user

import "context"

type DeleteUserRepository interface {
	DeleteUser(executionContext context.Context, userID string, user *User) error
}

type NotImplementedDeleteUserRepository struct{}

func (repository *NotImplementedDeleteUserRepository) DeleteUser(executionContext context.Context, userID string, user *User) error {
	return ErrNotImplemented
}
