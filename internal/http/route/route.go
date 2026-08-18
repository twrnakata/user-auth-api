package route

import (
	domainauth "backend-challenge-golang/internal/domain/auth"
	"backend-challenge-golang/internal/http/handler"

	"github.com/gofiber/fiber/v2"
)

type Route interface {
	InitRoute(application *fiber.App)
	InitRouteGroup(applicationGroup fiber.Router)
}

type route struct {
	AuthRegisterHandler *handler.AuthRegisterHandler
	AuthLoginHandler    *handler.AuthLoginHandler
}

func NewRoute(registerService domainauth.RegisterUserService, loginService domainauth.LoginUserService) Route {
	authRegisterHandler := &handler.AuthRegisterHandler{
		RegisterService: registerService,
	}
	authLoginHandler := &handler.AuthLoginHandler{
		LoginService: loginService,
	}

	return &route{
		AuthRegisterHandler: authRegisterHandler,
		AuthLoginHandler:    authLoginHandler,
	}
}

func NewApp(registerService domainauth.RegisterUserService, loginService domainauth.LoginUserService) *fiber.App {
	application := fiber.New()
	newRoute := NewRoute(registerService, loginService)
	newRoute.InitRoute(application)
	return application
}

func (route *route) InitRoute(application *fiber.App) {
	applicationGroup := application.Group("")
	route.InitRouteGroup(applicationGroup)
}

func (route *route) InitRouteGroup(applicationGroup fiber.Router) {
	route.healthGroup(applicationGroup)

	authGroup := applicationGroup.Group("/auth")
	route.authGroup(authGroup)
}

func (route *route) healthGroup(applicationGroup fiber.Router) {
	applicationGroup.Get("/health", func(context *fiber.Ctx) error {
		return context.Status(200).JSON(fiber.Map{
			"status": "ok",
		})
	})
}
