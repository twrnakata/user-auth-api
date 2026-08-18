package user

import (
	"context"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
)

type CountUserService struct {
	Repository domainuser.CountUserRepository
}

func (service *CountUserService) Count(executionContext context.Context, count *int64) error {
	if service.Repository == nil {
		return apperror.ErrCountUserRepositoryNotConfigured
	}
	if count == nil {
		return apperror.ErrCountResponseNil
	}

	return service.Repository.CountUsers(executionContext, count)
}
