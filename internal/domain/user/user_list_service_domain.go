package user

import "context"

type ListUserService interface {
	List(executionContext context.Context, users *[]User) error
}

type NotImplementedListUserService struct{}

func (service *NotImplementedListUserService) List(executionContext context.Context, users *[]User) error {
	return ErrNotImplemented
}
