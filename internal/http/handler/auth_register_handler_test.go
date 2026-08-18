package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	servicemodel "github.com/twrnakata/user-auth-api/internal/service/auth/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	domainauth "github.com/twrnakata/user-auth-api/internal/domain/auth"
	"github.com/twrnakata/user-auth-api/pkg/caller"
)

func requireHiddenInternalError(t *testing.T, response *http.Response, leakedMessage string) {
	t.Helper()
	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	responseBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeInternalError, int(responseEnvelope["code"].(float64)))
	require.Equal(t, "internal server error", responseEnvelope["message"])
	require.Equal(t, "internal server error", responseEnvelope["errors"])
	require.NotContains(t, string(responseBodyBytes), leakedMessage)
}

type fakeRegisterService struct {
	called   bool
	gotReq   servicemodel.RegisterUserRequestModel
	response *servicemodel.RegisterUserResponseModel
	err      error
}

func (service *fakeRegisterService) Register(executionContext context.Context, request *servicemodel.RegisterUserRequestModel, response *servicemodel.RegisterUserResponseModel) error {
	service.called = true
	if request != nil {
		service.gotReq = *request
	}
	if response != nil && service.response != nil {
		*response = *service.response
	}
	return service.err
}

func TestAuthRegisterHandler_InvalidJSON_Returns400(t *testing.T) {
	fakeService := &fakeRegisterService{}
	handler := &AuthRegisterHandler{RegisterService: fakeService}

	application := fiber.New()
	application.Post("/auth/register", handler.Register)

	request := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString(`{"name":`))
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

func TestAuthRegisterHandler_ValidBody_Returns201AndData(t *testing.T) {
	fakeService := &fakeRegisterService{
		response: &servicemodel.RegisterUserResponseModel{
			ID:        "u-123",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	handler := &AuthRegisterHandler{RegisterService: fakeService}

	application := fiber.New()
	application.Post("/auth/register", handler.Register)

	body := `{"name":" Alice ","email":"alice@example.com","password":" secret "}`
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

	responseData := responseEnvelope["data"].(map[string]any)
	require.Equal(t, "u-123", responseData["id"])
	require.Equal(t, "Alice", responseData["name"])
	require.Equal(t, "alice@example.com", responseData["email"])
	require.Equal(t, "2026-01-02 10:04:05", responseData["createdAt"])

	require.True(t, fakeService.called)
	require.Equal(t, "Alice", fakeService.gotReq.Name)
	require.Equal(t, "alice@example.com", fakeService.gotReq.Email)
	require.Equal(t, "secret", fakeService.gotReq.Password)
}

func TestAuthRegisterHandler_EmailAlreadyExists_Returns409(t *testing.T) {
	fakeService := &fakeRegisterService{
		err: domainauth.ErrEmailAlreadyExists,
	}
	handler := &AuthRegisterHandler{RegisterService: fakeService}

	application := fiber.New()
	application.Post("/auth/register", handler.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret"}`
	request := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusConflict, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeConflict, int(responseEnvelope["code"].(float64)))
}

func TestAuthRegisterHandler_ServiceError_Returns500WithoutLeaking(t *testing.T) {
	leakedMessage := "server selection timeout"
	fakeService := &fakeRegisterService{
		err: errors.New(leakedMessage),
	}
	handler := &AuthRegisterHandler{RegisterService: fakeService}

	application := fiber.New()
	application.Post("/auth/register", handler.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret"}`
	request := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	requireHiddenInternalError(t, response, leakedMessage)
}
