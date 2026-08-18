package httpapi

import (
	"backend-challenge-golang/internal/application/user"
	"backend-challenge-golang/internal/interfaces/http/handler"

	"github.com/gofiber/fiber/v2"
)

// NewApp wires HTTP routes for the challenge.
// This helper exists so we can smoke-test routing/middleware in a TDD workflow.
func NewApp(registerSvc user.RegisterUserService) *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		// health is intentionally lightweight; DB connectivity will come later.
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
		})
	})

	h := &handler.AuthRegisterHandler{
		RegisterService: registerSvc,
	}
	app.Post("/auth/register", h.Register)

	return app
}
