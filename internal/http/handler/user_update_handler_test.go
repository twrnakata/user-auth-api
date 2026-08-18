package handler

import (
	"bytes"
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

type fakeUpdateUserService struct {
	called  bool
	gotID   string
	gotUser domainuser.User
	user    *domainuser.User
	err     error
}

func (service *fakeUpdateUserService) Update(executionContext context.Context, userID string, user *domainuser.User) error {
	service.called = true
	service.gotID = userID
	if user != nil {
		service.gotUser = *user
	}
	if user != nil && service.user != nil {
		*user = *service.user
	}
	return service.err
}

func TestUserUpdateHandler_NameOnly_Returns200AndData(t *testing.T) {
	fakeService := &fakeUpdateUserService{
		user: &domainuser.User{
			ID:        "u-123",
			Name:      "Bob",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	handler := &UserUpdateHandler{UpdateUserService: fakeService}

	application := fiber.New()
	application.Put("/users/:id", handler.Update)

	body := `{"name":" Bob "}`
	request := httptest.NewRequest("PUT", "/users/u-123", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeSuccess, int(responseEnvelope["code"].(float64)))

	responseData := responseEnvelope["data"].(map[string]any)
	require.Equal(t, "u-123", responseData["id"])
	require.Equal(t, "Bob", responseData["name"])
	require.Equal(t, "alice@example.com", responseData["email"])
	require.Equal(t, "2026-01-02 10:04:05", responseData["createdAt"])
	require.True(t, fakeService.called)
	require.Equal(t, "u-123", fakeService.gotID)
	require.Equal(t, "Bob", fakeService.gotUser.Name)
	require.Equal(t, "", fakeService.gotUser.Email)
}

func TestUserUpdateHandler_EmailOnly_Returns200AndData(t *testing.T) {
	fakeService := &fakeUpdateUserService{
		user: &domainuser.User{
			ID:        "u-123",
			Name:      "Alice",
			Email:     "bob@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	handler := &UserUpdateHandler{UpdateUserService: fakeService}

	application := fiber.New()
	application.Put("/users/:id", handler.Update)

	body := `{"email":" bob@example.com "}`
	request := httptest.NewRequest("PUT", "/users/u-123", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)
	require.True(t, fakeService.called)
	require.Equal(t, "", fakeService.gotUser.Name)
	require.Equal(t, "bob@example.com", fakeService.gotUser.Email)
}

func TestUserUpdateHandler_EmptyBody_Returns400(t *testing.T) {
	fakeService := &fakeUpdateUserService{}
	handler := &UserUpdateHandler{UpdateUserService: fakeService}

	application := fiber.New()
	application.Put("/users/:id", handler.Update)

	request := httptest.NewRequest("PUT", "/users/u-123", bytes.NewBufferString(`{}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	require.False(t, fakeService.called)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeInvalidParam, int(responseEnvelope["code"].(float64)))
	require.Equal(t, "name or email is required", responseEnvelope["errors"])
}

func TestUserUpdateHandler_InvalidJSON_Returns400(t *testing.T) {
	fakeService := &fakeUpdateUserService{}
	handler := &UserUpdateHandler{UpdateUserService: fakeService}

	application := fiber.New()
	application.Put("/users/:id", handler.Update)

	request := httptest.NewRequest("PUT", "/users/u-123", bytes.NewBufferString(`{"name":`))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	require.False(t, fakeService.called)
}

func TestUserUpdateHandler_NotFound_Returns404(t *testing.T) {
	fakeService := &fakeUpdateUserService{
		err: domainuser.ErrUserNotFound,
	}
	handler := &UserUpdateHandler{UpdateUserService: fakeService}

	application := fiber.New()
	application.Put("/users/:id", handler.Update)

	body := `{"name":"Bob"}`
	request := httptest.NewRequest("PUT", "/users/missing", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)
}

func TestUserUpdateHandler_InvalidUserID_Returns404(t *testing.T) {
	fakeService := &fakeUpdateUserService{
		err: domainuser.ErrInvalidUserID,
	}
	handler := &UserUpdateHandler{UpdateUserService: fakeService}

	application := fiber.New()
	application.Put("/users/:id", handler.Update)

	body := `{"name":"Bob"}`
	request := httptest.NewRequest("PUT", "/users/not-an-object-id", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeNotFound, int(responseEnvelope["code"].(float64)))
	require.Equal(t, domainuser.ErrUserNotFound.Error(), responseEnvelope["errors"])
}

func TestUserUpdateHandler_EmailAlreadyExists_Returns409(t *testing.T) {
	fakeService := &fakeUpdateUserService{
		err: domainuser.ErrEmailAlreadyExists,
	}
	handler := &UserUpdateHandler{UpdateUserService: fakeService}

	application := fiber.New()
	application.Put("/users/:id", handler.Update)

	body := `{"email":"taken@example.com"}`
	request := httptest.NewRequest("PUT", "/users/u-123", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusConflict, response.StatusCode)

	responseBodyBytes, _ := io.ReadAll(response.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeConflict, int(responseEnvelope["code"].(float64)))
}
