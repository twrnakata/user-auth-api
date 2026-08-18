package route

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	servicemodel "backend-challenge-golang/internal/service/auth/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"backend-challenge-golang/pkg/caller"
)

type fakeRegisterService struct {
	response *servicemodel.RegisterUserResponse
}

func (service *fakeRegisterService) Register(executionContext context.Context, request *servicemodel.RegisterUserRequest, response *servicemodel.RegisterUserResponse) error {
	if response != nil && service.response != nil {
		*response = *service.response
	}
	return nil
}

type fakeLoginService struct {
	response *servicemodel.LoginUserResponse
}

func (service *fakeLoginService) Login(executionContext context.Context, request *servicemodel.LoginUserRequest, response *servicemodel.LoginUserResponse) error {
	if response != nil && service.response != nil {
		*response = *service.response
	}
	return nil
}

func TestNewApp_RoutesAreWired(t *testing.T) {
	application := NewApp(&fakeRegisterService{
		response: &servicemodel.RegisterUserResponse{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}, &fakeLoginService{
		response: &servicemodel.LoginUserResponse{
			Token: "jwt-token",
			ID:    "u-1",
			Name:  "Alice",
		},
	})

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
}
