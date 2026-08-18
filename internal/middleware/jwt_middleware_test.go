package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"backend-challenge-golang/pkg/caller"
	jwtpkg "backend-challenge-golang/pkg/jwt"
)

func TestJWT_MissingAuthorizationHeader_Returns401(t *testing.T) {
	jwtService, err := jwtpkg.NewJWTService("test-secret", jwtpkg.DefaultExpireDuration)
	require.NoError(t, err)

	application := fiber.New()
	application.Get("/users", JWT(jwtService), func(fiberContext *fiber.Ctx) error {
		return caller.Success(fiberContext, nil)
	})

	request := httptest.NewRequest("GET", "/users", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestJWT_InvalidBearerFormat_Returns401(t *testing.T) {
	jwtService, err := jwtpkg.NewJWTService("test-secret", jwtpkg.DefaultExpireDuration)
	require.NoError(t, err)

	application := fiber.New()
	application.Get("/users", JWT(jwtService), func(fiberContext *fiber.Ctx) error {
		return caller.Success(fiberContext, nil)
	})

	request := httptest.NewRequest("GET", "/users", nil)
	request.Header.Set("Authorization", "Token abc")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestJWT_InvalidToken_Returns401(t *testing.T) {
	jwtService, err := jwtpkg.NewJWTService("test-secret", jwtpkg.DefaultExpireDuration)
	require.NoError(t, err)

	application := fiber.New()
	application.Get("/users", JWT(jwtService), func(fiberContext *fiber.Ctx) error {
		return caller.Success(fiberContext, nil)
	})

	request := httptest.NewRequest("GET", "/users", nil)
	request.Header.Set("Authorization", "Bearer invalid-token")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}

func TestJWT_ValidToken_SetsLocalsAndCallsNext(t *testing.T) {
	jwtService, err := jwtpkg.NewJWTService("test-secret", jwtpkg.DefaultExpireDuration)
	require.NoError(t, err)

	token, err := jwtService.CreateToken("u-1", "Alice")
	require.NoError(t, err)

	application := fiber.New()
	application.Get("/users", JWT(jwtService), func(fiberContext *fiber.Ctx) error {
		require.Equal(t, "u-1", UserID(fiberContext))
		require.Equal(t, "Alice", UserName(fiberContext))
		return caller.Success(fiberContext, nil)
	})

	request := httptest.NewRequest("GET", "/users", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)
}
