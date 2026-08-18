package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	domainuser "backend-challenge-golang-7solution/internal/domain/user"
	"backend-challenge-golang-7solution/pkg/caller"
)

type fakeDeleteUserService struct {
	called bool
	gotID  string
	user   *domainuser.User
	err    error
}

func (service *fakeDeleteUserService) Delete(executionContext context.Context, userID string, user *domainuser.User) error {
	service.called = true
	service.gotID = userID
	if user != nil && service.user != nil {
		*user = *service.user
	}
	return service.err
}

func TestUserDeleteHandler_ValidID_Returns200AndDeletedUser(t *testing.T) {
	fakeService := &fakeDeleteUserService{
		user: &domainuser.User{
			ID:        "u-123",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	handler := &UserDeleteHandler{DeleteUserService: fakeService}

	application := fiber.New()
	application.Delete("/users/:id", handler.Delete)

	request := httptest.NewRequest("DELETE", "/users/u-123", nil)
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

func TestUserDeleteHandler_NotFound_Returns404(t *testing.T) {
	fakeService := &fakeDeleteUserService{
		err: domainuser.ErrUserNotFound,
	}
	handler := &UserDeleteHandler{DeleteUserService: fakeService}

	application := fiber.New()
	application.Delete("/users/:id", handler.Delete)

	request := httptest.NewRequest("DELETE", "/users/missing", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeNotFound, int(responseEnvelope["code"].(float64)))
}

func TestUserDeleteHandler_InvalidUserID_Returns404(t *testing.T) {
	fakeService := &fakeDeleteUserService{
		err: domainuser.ErrInvalidUserID,
	}
	handler := &UserDeleteHandler{DeleteUserService: fakeService}

	application := fiber.New()
	application.Delete("/users/:id", handler.Delete)

	request := httptest.NewRequest("DELETE", "/users/not-an-object-id", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeNotFound, int(responseEnvelope["code"].(float64)))
	require.Equal(t, domainuser.ErrUserNotFound.Error(), responseEnvelope["errors"])
}
