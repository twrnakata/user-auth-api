package handler

import (
	"context"

	handlermodel "github.com/twrnakata/user-auth-api/internal/http/handler/model"
	"github.com/gofiber/fiber/v2"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
	"github.com/twrnakata/user-auth-api/pkg/caller"
	"github.com/twrnakata/user-auth-api/pkg/datetime"
)

type UserListHandler struct {
	ListUserService domainuser.ListUserService
}

func (handler *UserListHandler) List(fiberContext *fiber.Ctx) error {
	if handler.ListUserService == nil {
		return caller.InternalError(fiberContext, apperror.ErrListUserServiceNotInitialized)
	}

	var users []domainuser.User
	err := handler.ListUserService.List(context.Background(), &users)
	if err != nil {
		return caller.InternalError(fiberContext, err)
	}

	responseUsers := make([]handlermodel.UserListItemModel, 0, len(users))
	for _, user := range users {
		responseUsers = append(responseUsers, handlermodel.UserListItemModel{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: datetime.FormatDateTime(user.CreatedAt),
		})
	}

	return caller.Success(fiberContext, handlermodel.UserListResponseModel{
		Users: responseUsers,
	})
}
