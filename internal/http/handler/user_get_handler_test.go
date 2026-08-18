package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	servicemodel "backend-challenge-golang/internal/service/user/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	domainuser "backend-challenge-golang/internal/domain/user"
	"backend-challenge-golang/pkg/caller"
)

type fakeGetUserService struct {
	called   bool
	gotID    string
	response *servicemodel.GetUserResponseModel
	err      error
}

func (service *fakeGetUserService) GetByID(executionContext context.Context, request *servicemodel.GetUserRequestModel, response *servicemodel.GetUserResponseModel) error {
	service.called = true
	if request != nil {
		service.gotID = request.ID
	}
	if response != nil && service.response != nil {
		*response = *service.response
	}
	return service.err
}

func TestUserGetHandler_ValidID_Returns200AndData(t *testing.T) {
	fakeService := &fakeGetUserService{
		response: &servicemodel.GetUserResponseModel{
			ID:        "u-123",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	handler := &UserGetHandler{GetUserService: fakeService}

	application := fiber.New()
	application.Get("/users/:id", handler.GetByID)

	request := httptest.NewRequest("GET", "/users/u-123", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeSuccess, int(responseEnvelope["code"].(float64)))

	responseData := responseEnvelope["data"].(map[string]any)
	require.Equal(t, "u-123", responseData["id"])
	require.Equal(t, "Alice", responseData["name"])
	require.Equal(t, "alice@example.com", responseData["email"])
	require.Equal(t, "2026-01-02 10:04:05", responseData["createdAt"])
	require.True(t, fakeService.called)
	require.Equal(t, "u-123", fakeService.gotID)
}

func TestUserGetHandler_NotFound_Returns404(t *testing.T) {
	fakeService := &fakeGetUserService{
		err: domainuser.ErrUserNotFound,
	}
	handler := &UserGetHandler{GetUserService: fakeService}

	application := fiber.New()
	application.Get("/users/:id", handler.GetByID)

	request := httptest.NewRequest("GET", "/users/missing", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeNotFound, int(responseEnvelope["code"].(float64)))
}

func TestUserGetHandler_InvalidUserID_Returns404(t *testing.T) {
	fakeService := &fakeGetUserService{
		err: domainuser.ErrInvalidUserID,
	}
	handler := &UserGetHandler{GetUserService: fakeService}

	application := fiber.New()
	application.Get("/users/:id", handler.GetByID)

	request := httptest.NewRequest("GET", "/users/not-an-object-id", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeNotFound, int(responseEnvelope["code"].(float64)))
	require.Equal(t, domainuser.ErrUserNotFound.Error(), responseEnvelope["errors"])
}
