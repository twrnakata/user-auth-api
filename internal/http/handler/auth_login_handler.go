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
)

// AuthLoginHandler handles POST /auth/login.
// It is an HTTP adapter that validates input and delegates to the domain/application port.
type AuthLoginHandler struct {
	LoginService domainauth.LoginUserService
}

func (handler *AuthLoginHandler) Login(c *fiber.Ctx) error {
	if handler.LoginService == nil {
		return caller.InternalServerError(c, "login service not initialized")
	}

	var request handlermodel.AuthLoginRequest
	if err := validateLoginRequest(c, &request); err != nil {
		return caller.BadRequest(c, err.Error())
	}

	var response servicemodel.LoginUserResponse
	err := handler.LoginService.Login(context.Background(), &servicemodel.LoginUserRequest{
		Email:    request.Email,
		Password: request.Password,
	}, &response)
	if err != nil {
		switch {
		case errors.Is(err, domainauth.ErrInvalidCredentials):
			return caller.Unauthorized(c, err.Error())
		case errors.Is(err, domainauth.ErrUserNotFound):
			return caller.NotFound(c, err.Error())
		default:
			return caller.InternalServerError(c, err.Error())
		}
	}

	responseBody := handlermodel.AuthLoginResponse{
		Token: response.Token,
		User: handlermodel.AuthLoginUserResponse{
			ID:   response.ID,
			Name: response.Name,
		},
	}

	return caller.Success(c, responseBody)
}

func validateLoginRequest(c *fiber.Ctx, request *handlermodel.AuthLoginRequest) error {
	if err := c.BodyParser(request); err != nil {
		return errors.New("invalid json body")
	}
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)
	if request.Email == "" || request.Password == "" {
		return errors.New("email and password are required")
	}
	return nil
}
