package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
	"github.com/twrnakata/user-auth-api/pkg/caller"
)

const validGetUserObjectID = "507f1f77bcf86cd799439011"

type fakeGetUserService struct {
	called bool
	gotID  string
	user   *domainuser.User
	err    error
}

func (service *fakeGetUserService) GetByID(executionContext context.Context, userID string, user *domainuser.User) error {
	service.called = true
	service.gotID = userID
	if user != nil && service.user != nil {
		*user = *service.user
	}
	return service.err
}

func TestUserGetHandler_ValidID_Returns200AndData(t *testing.T) {
	fakeService := &fakeGetUserService{
		user: &domainuser.User{
			ID:        validGetUserObjectID,
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	handler := &UserGetHandler{GetUserService: fakeService}

	application := fiber.New()
	application.Get("/users/:id", handler.GetByID)

	request := httptest.NewRequest("GET", "/users/"+validGetUserObjectID, nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeSuccess, int(responseEnvelope["code"].(float64)))

	responseData := responseEnvelope["data"].(map[string]any)
	require.Equal(t, validGetUserObjectID, responseData["id"])
	require.Equal(t, "Alice", responseData["name"])
	require.Equal(t, "alice@example.com", responseData["email"])
	require.Equal(t, "2026-01-02 10:04:05", responseData["createdAt"])
	require.True(t, fakeService.called)
	require.Equal(t, validGetUserObjectID, fakeService.gotID)
}

func TestUserGetHandler_NotFound_Returns404(t *testing.T) {
	fakeService := &fakeGetUserService{
		err: domainuser.ErrUserNotFound,
	}
	handler := &UserGetHandler{GetUserService: fakeService}

	application := fiber.New()
	application.Get("/users/:id", handler.GetByID)

	request := httptest.NewRequest("GET", "/users/"+validGetUserObjectID, nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)
	require.True(t, fakeService.called)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeNotFound, int(responseEnvelope["code"].(float64)))
	require.Equal(t, domainuser.ErrUserNotFound.Error(), responseEnvelope["errors"])
}

func TestUserGetHandler_InvalidUserID_Returns400WithoutCallingService(t *testing.T) {
	fakeService := &fakeGetUserService{}
	handler := &UserGetHandler{GetUserService: fakeService}

	application := fiber.New()
	application.Get("/users/:id", handler.GetByID)

	request := httptest.NewRequest("GET", "/users/not-an-object-id", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	require.False(t, fakeService.called)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeInvalidParam, int(responseEnvelope["code"].(float64)))
	require.Equal(t, apperror.ErrInvalidUserID.Error(), responseEnvelope["errors"])
}

func TestUserGetHandler_ServiceError_Returns500WithoutLeaking(t *testing.T) {
	leakedMessage := "server selection timeout"
	fakeService := &fakeGetUserService{
		err: errors.New(leakedMessage),
	}
	handler := &UserGetHandler{GetUserService: fakeService}

	application := fiber.New()
	application.Get("/users/:id", handler.GetByID)

	request := httptest.NewRequest("GET", "/users/"+validGetUserObjectID, nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	requireHiddenInternalError(t, response, leakedMessage)
}
