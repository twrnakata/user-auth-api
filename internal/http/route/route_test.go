package route

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	servicemodel "github.com/twrnakata/user-auth-api/internal/service/auth/model"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/caller"
	jwtpkg "github.com/twrnakata/user-auth-api/pkg/jwt"
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

type fakeListUserService struct {
	users []domainuser.User
}

func (service *fakeListUserService) List(executionContext context.Context, users *[]domainuser.User) error {
	if users != nil && service.users != nil {
		*users = service.users
	}
	return nil
}

type fakeGetUserService struct {
	user *domainuser.User
}

func (service *fakeGetUserService) GetByID(executionContext context.Context, userID string, user *domainuser.User) error {
	if user != nil && service.user != nil {
		*user = *service.user
	}
	return nil
}

type fakeUpdateUserService struct {
	user *domainuser.User
}

func (service *fakeUpdateUserService) Update(executionContext context.Context, userID string, user *domainuser.User) error {
	if user != nil && service.user != nil {
		*user = *service.user
	}
	return nil
}

type fakeDeleteUserService struct {
	user *domainuser.User
}

func (service *fakeDeleteUserService) Delete(executionContext context.Context, userID string, user *domainuser.User) error {
	if user != nil && service.user != nil {
		*user = *service.user
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
	}, &fakeListUserService{
		users: []domainuser.User{
			{
				ID:        "u-1",
				Name:      "Alice",
				Email:     "alice@example.com",
				CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
			},
		},
	}, &fakeGetUserService{
		user: &domainuser.User{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}, &fakeUpdateUserService{
		user: &domainuser.User{
			ID:        "u-1",
			Name:      "Bob",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}, &fakeDeleteUserService{
		user: &domainuser.User{
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

	getUserRequest := httptest.NewRequest("GET", "/users/507f1f77bcf86cd799439011", nil)
	getUserRequest.Header.Set("Authorization", "Bearer "+token)
	getUserResponse, err := application.Test(getUserRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, getUserResponse.StatusCode)

	listUsersRequest := httptest.NewRequest("GET", "/users", nil)
	listUsersRequest.Header.Set("Authorization", "Bearer "+token)
	listUsersResponse, err := application.Test(listUsersRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, listUsersResponse.StatusCode)

	updateUserBody := `{"name":"Bob"}`
	updateUserRequest := httptest.NewRequest("PUT", "/users/507f1f77bcf86cd799439011", bytes.NewBufferString(updateUserBody))
	updateUserRequest.Header.Set("Authorization", "Bearer "+token)
	updateUserRequest.Header.Set("Content-Type", "application/json")
	updateUserResponse, err := application.Test(updateUserRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, updateUserResponse.StatusCode)

	deleteUserRequest := httptest.NewRequest("DELETE", "/users/507f1f77bcf86cd799439011", nil)
	deleteUserRequest.Header.Set("Authorization", "Bearer "+token)
	deleteUserResponse, err := application.Test(deleteUserRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, deleteUserResponse.StatusCode)
}
