package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	domainauth "backend-challenge-golang/internal/domain/auth"
	"backend-challenge-golang/pkg/caller"
)

type fakeLoginService struct {
	called   bool
	response *servicemodel.LoginUserResponseModel
	err      error
}

func (service *fakeLoginService) Login(executionContext context.Context, request *servicemodel.LoginUserRequestModel, response *servicemodel.LoginUserResponseModel) error {
	service.called = true
	if response != nil && service.response != nil {
		*response = *service.response
	}
	return service.err
}

func TestAuthLoginHandler_InvalidJSON_Returns400(t *testing.T) {
	fakeService := &fakeLoginService{}
	handler := &AuthLoginHandler{LoginService: fakeService}

	application := fiber.New()
	application.Post("/auth/login", handler.Login)

	request := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(`{"email":`))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)

	require.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))

	require.Equal(t, caller.CodeInvalidParam, int(responseEnvelope["code"].(float64)))
	require.Equal(t, "invalid parameter", responseEnvelope["message"])
	require.False(t, fakeService.called)
}

func TestAuthLoginHandler_ValidBody_Returns200AndToken(t *testing.T) {
	fakeService := &fakeLoginService{
		response: &servicemodel.LoginUserResponseModel{
			Token: "jwt-token",
			ID:    "u-1",
			Name:  "Alice",
		},
	}
	handler := &AuthLoginHandler{LoginService: fakeService}

	application := fiber.New()
	application.Post("/auth/login", handler.Login)

	reqBody := `{"email":"alice@example.com","password":"secret"}`
	request := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)

	require.Equal(t, fiber.StatusOK, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))

	require.Equal(t, caller.CodeSuccess, int(responseEnvelope["code"].(float64)))
	require.Equal(t, "success", responseEnvelope["message"])

	responseData := responseEnvelope["data"].(map[string]any)
	require.Equal(t, "jwt-token", responseData["token"])

	responseUser := responseData["user"].(map[string]any)
	require.Equal(t, "u-1", responseUser["id"])
	require.Equal(t, "Alice", responseUser["name"])

	require.True(t, fakeService.called)
}

func TestAuthLoginHandler_InvalidCredentials_Returns401(t *testing.T) {
	fakeService := &fakeLoginService{
		err: domainauth.ErrInvalidCredentials,
	}
	handler := &AuthLoginHandler{LoginService: fakeService}

	application := fiber.New()
	application.Post("/auth/login", handler.Login)

	reqBody := `{"email":"alice@example.com","password":"wrong"}`
	request := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	request.Header.Set("Content-Type", "application/json")

	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusUnauthorized, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeUnauthorized, int(responseEnvelope["code"].(float64)))
}

func TestAuthLoginHandler_UserNotFound_Returns404(t *testing.T) {
	fakeService := &fakeLoginService{
		err: domainauth.ErrUserNotFound,
	}
	handler := &AuthLoginHandler{LoginService: fakeService}

	application := fiber.New()
	application.Post("/auth/login", handler.Login)

	reqBody := `{"email":"unknown@example.com","password":"secret"}`
	request := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	request.Header.Set("Content-Type", "application/json")

	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeNotFound, int(responseEnvelope["code"].(float64)))
}
