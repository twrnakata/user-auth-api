package route

import "github.com/gofiber/fiber/v2"

func (route *route) usersGroup(appGroup fiber.Router) {
	appGroup.Get("/", route.UserListHandler.List)
	appGroup.Get("/:id", route.UserGetHandler.GetByID)
	appGroup.Put("/:id", route.UserUpdateHandler.Update)
	appGroup.Delete("/:id", route.UserDeleteHandler.Delete)
}
