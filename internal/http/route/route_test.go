package route

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	servicemodel "backend-challenge-golang-7solution/internal/service/auth/model"
	userservicemodel "backend-challenge-golang-7solution/internal/service/user/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"backend-challenge-golang-7solution/pkg/caller"
	jwtpkg "backend-challenge-golang-7solution/pkg/jwt"
)

type fakeRegisterService struct {
	response *servicemodel.RegisterUserResponseModel
}

func (service *fakeRegisterService) Register(executionContext context.Context, request *servicemodel.RegisterUserRequestModel, response *servicemodel.RegisterUserResponseModel) error {
	if response != nil && service.response != nil {
		*response = *service.response
	}
	return nil
}

type fakeLoginService struct {
	response *servicemodel.LoginUserResponseModel
}

func (service *fakeLoginService) Login(executionContext context.Context, request *servicemodel.LoginUserRequestModel, response *servicemodel.LoginUserResponseModel) error {
	if response != nil && service.response != nil {
		*response = *service.response
	}
	return nil
}

type fakeGetUserService struct {
	response *userservicemodel.GetUserResponseModel
}

func (service *fakeGetUserService) GetByID(executionContext context.Context, request *userservicemodel.GetUserRequestModel, response *userservicemodel.GetUserResponseModel) error {
	if response != nil && service.response != nil {
		*response = *service.response
	}
	return nil
}

func TestNewApp_RoutesAreWired(t *testing.T) {
	jwtService, err := jwtpkg.NewJWTService("test-secret", jwtpkg.DefaultExpireDuration)
	require.NoError(t, err)

	application := NewApp(&fakeRegisterService{
		response: &servicemodel.RegisterUserResponseModel{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}, &fakeLoginService{
		response: &servicemodel.LoginUserResponseModel{
			Token: "jwt-token",
			ID:    "u-1",
			Name:  "Alice",
		},
	}, &fakeGetUserService{
		response: &userservicemodel.GetUserResponseModel{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}, jwtService)

	healthRequest := httptest.NewRequest("GET", "/health", nil)
	healthResponse, err := application.Test(healthRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, healthResponse.StatusCode)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret"}`
	request := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeSuccess, int(responseEnvelope["code"].(float64)))
	require.Equal(t, "success", responseEnvelope["message"])

	loginBody := `{"email":"alice@example.com","password":"secret"}`
	loginRequest := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(loginBody))
	loginRequest.Header.Set("Content-Type", "application/json")
	loginResponse, err := application.Test(loginRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, loginResponse.StatusCode)

	loginResponseBodyBytes, _ := io.ReadAll(loginResponse.Body)
	var loginResponseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(loginResponseBodyBytes, &loginResponseEnvelope))
	require.Equal(t, caller.CodeSuccess, int(loginResponseEnvelope["code"].(float64)))
	require.Equal(t, "success", loginResponseEnvelope["message"])

	token, err := jwtService.CreateToken("u-1", "Alice")
	require.NoError(t, err)

	getUserRequest := httptest.NewRequest("GET", "/users/u-1", nil)
	getUserRequest.Header.Set("Authorization", "Bearer "+token)
	getUserResponse, err := application.Test(getUserRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, getUserResponse.StatusCode)
}
