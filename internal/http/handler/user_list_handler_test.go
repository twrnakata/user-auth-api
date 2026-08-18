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

type fakeListUserService struct {
	called bool
	users  []domainuser.User
	err    error
}

func (service *fakeListUserService) List(executionContext context.Context, users *[]domainuser.User) error {
	service.called = true
	if users != nil && service.users != nil {
		*users = service.users
	}
	return service.err
}

func TestUserListHandler_List_Returns200AndUsers(t *testing.T) {
	fakeService := &fakeListUserService{
		users: []domainuser.User{
			{
				ID:        "u-1",
				Name:      "Alice",
				Email:     "alice@example.com",
				CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
			},
		},
	}
	handler := &UserListHandler{ListUserService: fakeService}

	application := fiber.New()
	application.Get("/users", handler.List)

	request := httptest.NewRequest("GET", "/users", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeSuccess, int(responseEnvelope["code"].(float64)))

	responseData := responseEnvelope["data"].(map[string]any)
	users := responseData["users"].([]any)
	require.Len(t, users, 1)
	user := users[0].(map[string]any)
	require.Equal(t, "u-1", user["id"])
	require.Equal(t, "Alice", user["name"])
	require.Equal(t, "alice@example.com", user["email"])
	require.Equal(t, "2026-01-02 10:04:05", user["createdAt"])
	require.True(t, fakeService.called)
}

func TestUserListHandler_List_Empty_Returns200AndEmptyArray(t *testing.T) {
	fakeService := &fakeListUserService{
		users: []domainuser.User{},
	}
	handler := &UserListHandler{ListUserService: fakeService}

	application := fiber.New()
	application.Get("/users", handler.List)

	request := httptest.NewRequest("GET", "/users", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeSuccess, int(responseEnvelope["code"].(float64)))

	responseData := responseEnvelope["data"].(map[string]any)
	users := responseData["users"].([]any)
	require.Empty(t, users)
}

func TestUserListHandler_List_ServiceError_Returns500(t *testing.T) {
	fakeService := &fakeListUserService{
		err: domainuser.ErrNotImplemented,
	}
	handler := &UserListHandler{ListUserService: fakeService}

	application := fiber.New()
	application.Get("/users", handler.List)

	request := httptest.NewRequest("GET", "/users", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}
