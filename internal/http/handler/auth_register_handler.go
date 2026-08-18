// Package handler contains Fiber HTTP handlers (adapters).
package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	handlermodel "github.com/twrnakata/user-auth-api/internal/http/handler/model"
	servicemodel "github.com/twrnakata/user-auth-api/internal/service/auth/model"

	domainauth "github.com/twrnakata/user-auth-api/internal/domain/auth"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
	"github.com/twrnakata/user-auth-api/pkg/caller"
	"github.com/twrnakata/user-auth-api/pkg/datetime"
	"github.com/twrnakata/user-auth-api/pkg/validation"
)

type AuthRegisterHandler struct {
	RegisterService domainauth.RegisterUserService
}

// Register handles POST /auth/register.
func (handler *AuthRegisterHandler) Register(c *fiber.Ctx) error {
	if handler.RegisterService == nil {
		return caller.InternalError(c, apperror.ErrRegisterServiceNotInitialized)
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
		return caller.InternalError(c, err)
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
		return apperror.ErrInvalidJSONBody
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)
	if request.Name == "" || request.Email == "" || request.Password == "" {
		return apperror.ErrNameEmailPasswordRequired
	}
	if !validation.IsValidEmail(request.Email) {
		return apperror.ErrInvalidEmail
	}
	return nil
}
