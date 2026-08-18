package handler

import (
	"context"
	"errors"
	"strings"

	handlermodel "backend-challenge-golang-7solution/internal/http/handler/model"
	"github.com/gofiber/fiber/v2"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	"backend-challenge-golang-7solution/pkg/caller"
	"backend-challenge-golang-7solution/pkg/datetime"
)

type UserUpdateHandler struct {
	UpdateUserService domainuser.UpdateUserService
}

func (handler *UserUpdateHandler) Update(fiberContext *fiber.Ctx) error {
	if handler.UpdateUserService == nil {
		return caller.InternalServerError(fiberContext, "update user service not initialized")
	}

	var request handlermodel.UserUpdateRequestModel
	err := validateUpdateUserRequest(fiberContext, &request)
	if err != nil {
		return caller.BadRequest(fiberContext, err.Error())
	}

	user := domainuser.User{
		Name:  request.Name,
		Email: request.Email,
	}
	err = handler.UpdateUserService.Update(context.Background(), request.ID, &user)
	if err != nil {
		switch {
		case errors.Is(err, domainuser.ErrInvalidUserID):
			return caller.NotFound(fiberContext, domainuser.ErrUserNotFound.Error())
		case errors.Is(err, domainuser.ErrUserNotFound):
			return caller.NotFound(fiberContext, err.Error())
		case errors.Is(err, domainuser.ErrEmailAlreadyExists):
			return caller.Conflict(fiberContext, err.Error())
		default:
			return caller.InternalServerError(fiberContext, err.Error())
		}
	}

	responseBody := handlermodel.UserUpdateResponseModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: datetime.FormatDateTime(user.CreatedAt),
	}
	return caller.Success(fiberContext, responseBody)
}

func validateUpdateUserRequest(fiberContext *fiber.Ctx, request *handlermodel.UserUpdateRequestModel) error {
	if err := fiberContext.BodyParser(request); err != nil {
		return errors.New("invalid json body")
	}

	request.ID = strings.TrimSpace(fiberContext.Params("id"))
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	if request.ID == "" {
		return errors.New("id is required")
	}
	if request.Name == "" && request.Email == "" {
		return errors.New("name or email is required")
	}
	return nil
}
