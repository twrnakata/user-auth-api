package route

import "github.com/gofiber/fiber/v2"

func (route *route) authGroup(applicationGroup fiber.Router) {
	applicationGroup.Post("/register", route.AuthRegisterHandler.Register)
	applicationGroup.Post("/login", route.AuthLoginHandler.Login)
}
