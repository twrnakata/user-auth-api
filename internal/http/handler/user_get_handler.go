package handler

import (
	"context"
	"errors"
	"strings"

	handlermodel "backend-challenge-golang-7solution/internal/http/handler/model"
	servicemodel "backend-challenge-golang-7solution/internal/service/user/model"
	"github.com/gofiber/fiber/v2"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	"backend-challenge-golang-7solution/pkg/caller"
	"backend-challenge-golang-7solution/pkg/datetime"
)

type UserGetHandler struct {
	GetUserService domainuser.GetUserService
}

func (handler *UserGetHandler) GetByID(fiberContext *fiber.Ctx) error {
	if handler.GetUserService == nil {
		return caller.InternalServerError(fiberContext, "get user service not initialized")
	}

	var request handlermodel.UserGetRequestModel
	err := validateGetUserRequest(fiberContext, &request)
	if err != nil {
		return caller.BadRequest(fiberContext, err.Error())
	}

	var response servicemodel.GetUserResponseModel
	err = handler.GetUserService.GetByID(context.Background(), &servicemodel.GetUserRequestModel{
		ID: request.ID,
	}, &response)
	if err != nil {
		switch {
		case errors.Is(err, domainuser.ErrInvalidUserID):
			return caller.NotFound(fiberContext, domainuser.ErrUserNotFound.Error())
		case errors.Is(err, domainuser.ErrUserNotFound):
			return caller.NotFound(fiberContext, err.Error())
		default:
			return caller.InternalServerError(fiberContext, err.Error())
		}
	}

	responseBody := handlermodel.UserGetResponseModel{
		ID:        response.ID,
		Name:      response.Name,
		Email:     response.Email,
		CreatedAt: datetime.FormatDateTime(response.CreatedAt),
	}
	return caller.Success(fiberContext, responseBody)
}

func validateGetUserRequest(fiberContext *fiber.Ctx, request *handlermodel.UserGetRequestModel) error {
	request.ID = strings.TrimSpace(fiberContext.Params("id"))
	if request.ID == "" {
		return errors.New("id is required")
	}
	return nil
}
