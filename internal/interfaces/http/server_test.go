package httpapi

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

	"backend-challenge-golang/internal/application/user"
	"backend-challenge-golang/pkg/caller"
)

type fakeRegisterService struct {
	resp user.RegisterUserResponse
}

func (s *fakeRegisterService) Register(ctx context.Context, req user.RegisterUserRequest) (user.RegisterUserResponse, error) {
	return s.resp, nil
}

func TestNewApp_RoutesAreWired(t *testing.T) {
	app := NewApp(&fakeRegisterService{
		resp: user.RegisterUserResponse{
			ID:        "u-1",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	})

	// /health
	reqHealth := httptest.NewRequest("GET", "/health", nil)
	respHealth, err := app.Test(reqHealth, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, respHealth.StatusCode)

	// /auth/register
	body := `{"name":"Alice","email":"alice@example.com","password":"secret"}`
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)

	b, _ := io.ReadAll(resp.Body)
	var env map[string]any
	require.NoError(t, json.Unmarshal(b, &env))
	require.Equal(t, caller.CodeSuccess, int(env["code"].(float64)))
	require.Equal(t, "success", env["message"])
}

