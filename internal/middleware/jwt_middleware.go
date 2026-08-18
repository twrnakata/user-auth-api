package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/twrnakata/user-auth-api/pkg/apperror"
	"github.com/twrnakata/user-auth-api/pkg/caller"
	jwtpkg "github.com/twrnakata/user-auth-api/pkg/jwt"
)

const (
	LocalKeyUserID = "userId"
	LocalKeyName   = "name"
	bearerPrefix   = "Bearer "
)

func JWT(jwtService *jwtpkg.JWTService) fiber.Handler {
	return func(fiberContext *fiber.Ctx) error {
		if jwtService == nil {
			return caller.InternalError(fiberContext, apperror.ErrJWTServiceNotInitialized)
		}

		authorizationHeader := strings.TrimSpace(fiberContext.Get("Authorization"))
		if authorizationHeader == "" {
			return caller.Unauthorized(fiberContext, apperror.ErrMissingAuthorizationHeader.Error())
		}
		if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
			return caller.Unauthorized(fiberContext, apperror.ErrInvalidAuthorizationHeaderFormat.Error())
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, bearerPrefix))
		if tokenString == "" {
			return caller.Unauthorized(fiberContext, apperror.ErrMissingBearerToken.Error())
		}

		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			return caller.Unauthorized(fiberContext, apperror.ErrInvalidOrExpiredToken.Error())
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
