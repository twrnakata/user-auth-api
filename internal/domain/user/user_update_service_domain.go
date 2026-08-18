package user

import (
	"context"
	"errors"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type UpdateUserService interface {
	Update(executionContext context.Context, userID string, user *User) error
}

type NotImplementedUpdateUserService struct{}

func (service *NotImplementedUpdateUserService) Update(executionContext context.Context, userID string, user *User) error {
	return ErrNotImplemented
}
