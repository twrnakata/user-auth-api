// Package handler contains Fiber HTTP handlers (adapters).
package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	handlermodel "backend-challenge-golang/internal/http/handler/model"
	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"github.com/gofiber/fiber/v2"

	domainauth "backend-challenge-golang/internal/domain/auth"
	"backend-challenge-golang/pkg/caller"
)

type AuthRegisterHandler struct {
	RegisterService domainauth.RegisterUserService
}

// Register handles POST /auth/register.
func (handler *AuthRegisterHandler) Register(c *fiber.Ctx) error {
	if handler.RegisterService == nil {
		return caller.InternalServerError(c, "register service not initialized")
	}

	var request handlermodel.AuthRegisterRequest
	err := validateRegisterRequest(c, &request)
	if err != nil {
		return caller.BadRequest(c, err.Error())
	}

	var response servicemodel.RegisterUserResponse
	err = handler.RegisterService.Register(context.Background(), &servicemodel.RegisterUserRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}, &response)
	if err != nil {
		return caller.InternalServerError(c, err.Error())
	}

	responseBody := handlermodel.AuthRegisterResponse{
		ID:        response.ID,
		Name:      response.Name,
		Email:     response.Email,
		CreatedAt: response.CreatedAt.Format(time.RFC3339),
	}
	return caller.Created(c, responseBody)
}

func validateRegisterRequest(c *fiber.Ctx, request *handlermodel.AuthRegisterRequest) error {
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
