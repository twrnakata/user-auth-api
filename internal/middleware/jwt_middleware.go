package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"backend-challenge-golang-7solution/pkg/caller"
	jwtpkg "backend-challenge-golang-7solution/pkg/jwt"
)

const (
	LocalKeyUserID = "userId"
	LocalKeyName   = "name"
	bearerPrefix   = "Bearer "
)

func JWT(jwtService *jwtpkg.JWTService) fiber.Handler {
	return func(fiberContext *fiber.Ctx) error {
		if jwtService == nil {
			return caller.InternalServerError(fiberContext, "jwt service not initialized")
		}

		authorizationHeader := strings.TrimSpace(fiberContext.Get("Authorization"))
		if authorizationHeader == "" {
			return caller.Unauthorized(fiberContext, "missing authorization header")
		}
		if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
			return caller.Unauthorized(fiberContext, "invalid authorization header format")
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, bearerPrefix))
		if tokenString == "" {
			return caller.Unauthorized(fiberContext, "missing bearer token")
		}

		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			return caller.Unauthorized(fiberContext, "invalid or expired token")
		}

		fiberContext.Locals(LocalKeyUserID, claims.UserID)
		fiberContext.Locals(LocalKeyName, claims.Name)
		return fiberContext.Next()
	}
}

func UserID(fiberContext *fiber.Ctx) string {
	userID, _ := fiberContext.Locals(LocalKeyUserID).(string)
	return userID
}

func UserName(fiberContext *fiber.Ctx) string {
	name, _ := fiberContext.Locals(LocalKeyName).(string)
	return name
}
