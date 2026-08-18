package route

import (
	domainauth "backend-challenge-golang/internal/domain/auth"
	domainuser "backend-challenge-golang/internal/domain/user"
	"backend-challenge-golang/internal/http/handler"
	"backend-challenge-golang/internal/middleware"
	jwtpkg "backend-challenge-golang/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type Route interface {
	InitRoute(application *fiber.App)
	InitRouteGroup(applicationGroup fiber.Router)
}

type route struct {
	AuthRegisterHandler *handler.AuthRegisterHandler
	AuthLoginHandler    *handler.AuthLoginHandler
	UserGetHandler      *handler.UserGetHandler
	JWTService          *jwtpkg.JWTService
}

func NewRoute(registerService domainauth.RegisterUserService, loginService domainauth.LoginUserService, getUserService domainuser.GetUserService, jwtService *jwtpkg.JWTService) Route {
	authRegisterHandler := &handler.AuthRegisterHandler{
		RegisterService: registerService,
	}
	authLoginHandler := &handler.AuthLoginHandler{
		LoginService: loginService,
	}
	userGetHandler := &handler.UserGetHandler{
		GetUserService: getUserService,
	}

	return &route{
		AuthRegisterHandler: authRegisterHandler,
		AuthLoginHandler:    authLoginHandler,
		UserGetHandler:      userGetHandler,
		JWTService:          jwtService,
	}
}

func NewApp(registerService domainauth.RegisterUserService, loginService domainauth.LoginUserService, getUserService domainuser.GetUserService, jwtService *jwtpkg.JWTService) *fiber.App {
	application := fiber.New()
	newRoute := NewRoute(registerService, loginService, getUserService, jwtService)
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

	usersGroup := applicationGroup.Group("/users", middleware.JWT(route.JWTService))
	route.usersGroup(usersGroup)
}

func (route *route) healthGroup(appGroup fiber.Router) {
	appGroup.Get("/health", func(context *fiber.Ctx) error {
		return context.Status(200).JSON(fiber.Map{
			"status": "ok",
		})
	})
}
