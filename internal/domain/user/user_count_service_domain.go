package user

import "context"

type CountUserService interface {
	Count(executionContext context.Context, count *int64) error
}

type NotImplementedCountUserService struct{}

func (service *NotImplementedCountUserService) Count(executionContext context.Context, count *int64) error {
	return ErrNotImplemented
}
