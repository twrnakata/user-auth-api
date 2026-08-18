package handler

import (
	"context"
	"errors"
	"strings"

	handlermodel "github.com/twrnakata/user-auth-api/internal/http/handler/model"
	"github.com/gofiber/fiber/v2"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
	"github.com/twrnakata/user-auth-api/pkg/caller"
	"github.com/twrnakata/user-auth-api/pkg/datetime"
)

type UserDeleteHandler struct {
	DeleteUserService domainuser.DeleteUserService
}

func (handler *UserDeleteHandler) Delete(fiberContext *fiber.Ctx) error {
	if handler.DeleteUserService == nil {
		return caller.InternalError(fiberContext, apperror.ErrDeleteUserServiceNotInitialized)
	}

	var request handlermodel.UserDeleteRequestModel
	err := validateDeleteUserRequest(fiberContext, &request)
	if err != nil {
		return caller.BadRequest(fiberContext, err.Error())
	}

	var user domainuser.User
	err = handler.DeleteUserService.Delete(context.Background(), request.ID, &user)
	if err != nil {
		switch {
		case errors.Is(err, domainuser.ErrInvalidUserID):
			return caller.NotFound(fiberContext, domainuser.ErrUserNotFound.Error())
		case errors.Is(err, domainuser.ErrUserNotFound):
			return caller.NotFound(fiberContext, err.Error())
		default:
			return caller.InternalError(fiberContext, err)
		}
	}

	responseBody := handlermodel.UserDeleteResponseModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: datetime.FormatDateTime(user.CreatedAt),
	}
	return caller.Success(fiberContext, responseBody)
}

func validateDeleteUserRequest(fiberContext *fiber.Ctx, request *handlermodel.UserDeleteRequestModel) error {
	request.ID = strings.TrimSpace(fiberContext.Params("id"))
	if request.ID == "" {
		return apperror.ErrIDRequired
	}
	return nil
}
