package handler

import (
	"context"

	handlermodel "backend-challenge-golang-7solution/internal/http/handler/model"
	"github.com/gofiber/fiber/v2"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	"backend-challenge-golang-7solution/pkg/caller"
	"backend-challenge-golang-7solution/pkg/datetime"
)

type UserListHandler struct {
	ListUserService domainuser.ListUserService
}

func (handler *UserListHandler) List(fiberContext *fiber.Ctx) error {
	if handler.ListUserService == nil {
		return caller.InternalServerError(fiberContext, "list user service not initialized")
	}

	var users []domainuser.User
	err := handler.ListUserService.List(context.Background(), &users)
	if err != nil {
		return caller.InternalServerError(fiberContext, err.Error())
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
