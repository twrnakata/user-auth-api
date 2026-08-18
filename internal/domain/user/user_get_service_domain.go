package user

import (
	"context"
	"errors"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidUserID  = errors.New("invalid user id")
)

type GetUserService interface {
	GetByID(executionContext context.Context, userID string, user *User) error
}

type NotImplementedGetUserService struct{}

func (service *NotImplementedGetUserService) GetByID(executionContext context.Context, userID string, user *User) error {
	return ErrNotImplemented
}
