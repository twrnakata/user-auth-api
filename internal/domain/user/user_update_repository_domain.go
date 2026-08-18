package user

import "context"

type UpdateUserRepository interface {
	UpdateUser(executionContext context.Context, userID string, update UserUpdate, user *User) error
}

type NotImplementedUpdateUserRepository struct{}

func (repository *NotImplementedUpdateUserRepository) UpdateUser(executionContext context.Context, userID string, update UserUpdate, user *User) error {
	return ErrNotImplemented
}
