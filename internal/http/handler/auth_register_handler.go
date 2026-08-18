// Package handler contains Fiber HTTP handlers (adapters).
package handler

import (
	"context"
	"errors"
	"strings"

	handlermodel "backend-challenge-golang/internal/http/handler/model"
	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"github.com/gofiber/fiber/v2"

	domainauth "backend-challenge-golang/internal/domain/auth"
	"backend-challenge-golang/pkg/caller"
	"backend-challenge-golang/pkg/datetime"
)

type AuthRegisterHandler struct {
	RegisterService domainauth.RegisterUserService
}

// Register handles POST /auth/register.
func (handler *AuthRegisterHandler) Register(c *fiber.Ctx) error {
	if handler.RegisterService == nil {
		return caller.InternalServerError(c, "register service not initialized")
	}

	var request handlermodel.AuthRegisterRequestModel
	err := validateRegisterRequest(c, &request)
	if err != nil {
		return caller.BadRequest(c, err.Error())
	}

	var response servicemodel.RegisterUserResponseModel
	err = handler.RegisterService.Register(context.Background(), &servicemodel.RegisterUserRequestModel{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}, &response)
	if err != nil {
		if errors.Is(err, domainauth.ErrEmailAlreadyExists) {
			return caller.Conflict(c, err.Error())
		}
		return caller.InternalServerError(c, err.Error())
	}

	responseBody := handlermodel.AuthRegisterResponseModel{
		ID:        response.ID,
		Name:      response.Name,
		Email:     response.Email,
		CreatedAt: datetime.FormatDateTime(response.CreatedAt),
	}
	return caller.Created(c, responseBody)
}

func validateRegisterRequest(c *fiber.Ctx, request *handlermodel.AuthRegisterRequestModel) error {
	if err := c.BodyParser(request); err != nil {
		return errors.New("invalid json body")
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)
	if request.Name == "" || request.Email == "" || request.Password == "" {
		return errors.New("name, email, password are required")
	}
	return nil
}
