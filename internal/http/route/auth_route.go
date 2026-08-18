package route

import "github.com/gofiber/fiber/v2"

func (route *route) authGroup(appGroup fiber.Router) {
	appGroup.Post("/register", route.AuthRegisterHandler.Register)
	appGroup.Post("/login", route.AuthLoginHandler.Login)
}
