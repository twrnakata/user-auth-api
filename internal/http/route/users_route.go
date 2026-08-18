package route

import "github.com/gofiber/fiber/v2"

func (route *route) usersGroup(appGroup fiber.Router) {
	appGroup.Get("/:id", route.UserGetHandler.GetByID)
}
