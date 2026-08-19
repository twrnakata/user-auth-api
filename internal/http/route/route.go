package route

import (
	"time"

	domainauth "github.com/twrnakata/user-auth-api/internal/domain/auth"
	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/internal/http/handler"
	"github.com/twrnakata/user-auth-api/internal/middleware"
	jwtpkg "github.com/twrnakata/user-auth-api/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type Route interface {
	InitRoute(appGroup *fiber.App)
	InitRouteGroup(appGroup fiber.Router)
}

type route struct {
	AuthRegisterHandler *handler.AuthRegisterHandler
	AuthLoginHandler    *handler.AuthLoginHandler
	UserListHandler     *handler.UserListHandler
	UserGetHandler      *handler.UserGetHandler
	UserUpdateHandler   *handler.UserUpdateHandler
	UserDeleteHandler   *handler.UserDeleteHandler
	JWTService          *jwtpkg.JWTService
}

func NewRoute(registerService domainauth.RegisterUserService, loginService domainauth.LoginUserService, listUserService domainuser.ListUserService, getUserService domainuser.GetUserService, updateUserService domainuser.UpdateUserService, deleteUserService domainuser.DeleteUserService, jwtService *jwtpkg.JWTService) Route {
	authRegisterHandler := &handler.AuthRegisterHandler{
		RegisterService: registerService,
	}
	authLoginHandler := &handler.AuthLoginHandler{
		LoginService: loginService,
	}
	userListHandler := &handler.UserListHandler{
		ListUserService: listUserService,
	}
	userGetHandler := &handler.UserGetHandler{
		GetUserService: getUserService,
	}
	userUpdateHandler := &handler.UserUpdateHandler{
		UpdateUserService: updateUserService,
	}
	userDeleteHandler := &handler.UserDeleteHandler{
		DeleteUserService: deleteUserService,
	}

	return &route{
		AuthRegisterHandler: authRegisterHandler,
		AuthLoginHandler:    authLoginHandler,
		UserListHandler:     userListHandler,
		UserGetHandler:      userGetHandler,
		UserUpdateHandler:   userUpdateHandler,
		UserDeleteHandler:   userDeleteHandler,
		JWTService:          jwtService,
	}
}

const fiberReadTimeout = 30 * time.Second

func NewApp(registerService domainauth.RegisterUserService, loginService domainauth.LoginUserService, listUserService domainuser.ListUserService, getUserService domainuser.GetUserService, updateUserService domainuser.UpdateUserService, deleteUserService domainuser.DeleteUserService, jwtService *jwtpkg.JWTService) *fiber.App {
	application := fiber.New(fiber.Config{
		ReadTimeout: fiberReadTimeout,
	})
	newRoute := NewRoute(registerService, loginService, listUserService, getUserService, updateUserService, deleteUserService, jwtService)
	newRoute.InitRoute(application)
	return application
}

func (route *route) InitRoute(appGroup *fiber.App) {
	appGroup.Use(middleware.Logging(nil))
	appGroup.Use(middleware.Recover(nil))
	route.InitRouteGroup(appGroup.Group(""))
}

func (route *route) InitRouteGroup(appGroup fiber.Router) {
	route.healthGroup(appGroup)

	authGroup := appGroup.Group("/auth")
	route.authGroup(authGroup)

	usersGroup := appGroup.Group("/users", middleware.JWT(route.JWTService))
	route.usersGroup(usersGroup)
}

func (route *route) healthGroup(appGroup fiber.Router) {
	appGroup.Get("/health", func(context *fiber.Ctx) error {
		return context.Status(200).JSON(fiber.Map{
			"status": "ok",
		})
	})
}
