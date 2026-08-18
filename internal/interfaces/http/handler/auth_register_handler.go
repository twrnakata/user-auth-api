package handler

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend-challenge-golang/internal/application/user"
	"backend-challenge-golang/pkg/caller"
)

type AuthRegisterHandler struct {
	RegisterService user.RegisterUserService
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

// Register handles POST /auth/register.
func (h *AuthRegisterHandler) Register(c *fiber.Ctx) error {
	if h.RegisterService == nil {
		return caller.InternalServerError(c, "register service not initialized")
	}

	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return caller.BadRequest(c, "invalid json body")
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return caller.BadRequest(c, "name, email, password are required")
	}

	resp, err := h.RegisterService.Register(context.Background(), user.RegisterUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return caller.InternalServerError(c, err.Error())
	}

	out := registerResponse{
		ID:        resp.ID,
		Name:      resp.Name,
		Email:     resp.Email,
		CreatedAt: resp.CreatedAt.Format(time.RFC3339),
	}
	return caller.Created(c, out)
}

