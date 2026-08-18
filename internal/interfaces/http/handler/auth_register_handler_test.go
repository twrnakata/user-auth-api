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

	"backend-challenge-golang/internal/application/user"
	"backend-challenge-golang/pkg/caller"
)

type fakeRegisterService struct {
	called bool
	gotReq user.RegisterUserRequest
	resp   user.RegisterUserResponse
	err    error
}

func (s *fakeRegisterService) Register(ctx context.Context, req user.RegisterUserRequest) (user.RegisterUserResponse, error) {
	s.called = true
	s.gotReq = req
	return s.resp, s.err
}

func TestAuthRegisterHandler_InvalidJSON_Returns400(t *testing.T) {
	fakeSvc := &fakeRegisterService{}
	h := &AuthRegisterHandler{RegisterService: fakeSvc}

	app := fiber.New()
	app.Post("/auth/register", h.Register)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString(`{"name":`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	b, _ := io.ReadAll(resp.Body)
	var env map[string]any
	require.NoError(t, json.Unmarshal(b, &env))

	require.Equal(t, caller.CodeInvalidParam, int(env["code"].(float64)))
	require.Equal(t, "invalid parameter", env["message"])
	require.False(t, fakeSvc.called)
}

func TestAuthRegisterHandler_ValidBody_Returns201AndData(t *testing.T) {
	fakeSvc := &fakeRegisterService{
		resp: user.RegisterUserResponse{
			ID:        "u-123",
			Name:      "Alice",
			Email:     "alice@example.com",
			CreatedAt: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}
	h := &AuthRegisterHandler{RegisterService: fakeSvc}

	app := fiber.New()
	app.Post("/auth/register", h.Register)

	body := `{"name":" Alice ","email":"alice@example.com","password":" secret "}`
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

	data := env["data"].(map[string]any)
	require.Equal(t, "u-123", data["id"])
	require.Equal(t, "Alice", data["name"])
	require.Equal(t, "alice@example.com", data["email"])
	require.Equal(t, "2026-01-02T03:04:05Z", data["createdAt"])

	require.True(t, fakeSvc.called)
	require.Equal(t, "Alice", fakeSvc.gotReq.Name)
	require.Equal(t, "alice@example.com", fakeSvc.gotReq.Email)
	require.Equal(t, "secret", fakeSvc.gotReq.Password)
}

