package handler

import (
	"context"
	"errors"
	"strings"

	handlermodel "github.com/twrnakata/user-auth-api/internal/http/handler/model"
	servicemodel "github.com/twrnakata/user-auth-api/internal/service/auth/model"
	"github.com/gofiber/fiber/v2"

	domainauth "github.com/twrnakata/user-auth-api/internal/domain/auth"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
	"github.com/twrnakata/user-auth-api/pkg/caller"
)

// AuthLoginHandler handles POST /auth/login.
// It is an HTTP adapter that validates input and delegates to the domain/application port.
type AuthLoginHandler struct {
	LoginService domainauth.LoginUserService
}

func (handler *AuthLoginHandler) Login(c *fiber.Ctx) error {
	if handler.LoginService == nil {
		return caller.InternalError(c, apperror.ErrLoginServiceNotInitialized)
	}

	var request handlermodel.AuthLoginRequestModel
	if err := validateLoginRequest(c, &request); err != nil {
		return caller.BadRequest(c, err.Error())
	}

	var response servicemodel.LoginUserResponseModel
	err := handler.LoginService.Login(context.Background(), &servicemodel.LoginUserRequestModel{
		Email:    request.Email,
		Password: request.Password,
	}, &response)
	if err != nil {
		// Unknown email and wrong password both surface as ErrInvalidCredentials (401).
		// A 404 here would let callers enumerate registered emails.
		if errors.Is(err, domainauth.ErrInvalidCredentials) {
			return caller.Unauthorized(c, err.Error())
		}
		return caller.InternalError(c, err)
	}

	responseBody := handlermodel.AuthLoginResponseModel{
		Token: response.Token,
		User: handlermodel.AuthLoginUserResponseModel{
			ID:   response.ID,
			Name: response.Name,
		},
	}

	return caller.Success(c, responseBody)
}

func validateLoginRequest(c *fiber.Ctx, request *handlermodel.AuthLoginRequestModel) error {
	if err := c.BodyParser(request); err != nil {
		return apperror.ErrInvalidJSONBody
	}
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)
	if request.Email == "" || request.Password == "" {
		return apperror.ErrEmailAndPasswordRequired
	}
	return nil
}
