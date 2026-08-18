package caller

import (
	"net/http"
	"time"

	"backend-challenge-golang-7solution/pkg/datetime"
	"github.com/gofiber/fiber/v2"
)

// Numeric response codes (application-level code).
// Note: these are not HTTP status codes.
const (
	CodeSuccess       = 0
	CodeInvalidParam  = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeConflict      = 409
	CodeInternalError = 500
)

// ResponseModel is the standard JSON envelope.
type ResponseModel struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Errors     any    `json:"errors,omitempty"`
	ServerTime string `json:"serverTime"`
}

func now() string {
	return datetime.GetCurrentDateTimeNow().Format(time.RFC3339)
}

func Success(c *fiber.Ctx, data any) error {
	return c.Status(http.StatusOK).JSON(ResponseModel{
		Code:       CodeSuccess,
		Message:    "success",
		Data:       data,
		ServerTime: now(),
	})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(http.StatusCreated).JSON(ResponseModel{
		Code:       CodeSuccess,
		Message:    "success",
		Data:       data,
		ServerTime: now(),
	})
}

func BadRequest(c *fiber.Ctx, errs any) error {
	return c.Status(http.StatusBadRequest).JSON(ResponseModel{
		Code:       CodeInvalidParam,
		Message:    "invalid parameter",
		Errors:     errs,
		ServerTime: now(),
	})
}

func Unauthorized(c *fiber.Ctx, errs any) error {
	return c.Status(http.StatusUnauthorized).JSON(ResponseModel{
		Code:       CodeUnauthorized,
		Message:    "unauthorized",
		Errors:     errs,
		ServerTime: now(),
	})
}

func Forbidden(c *fiber.Ctx, errs any) error {
	return c.Status(http.StatusForbidden).JSON(ResponseModel{
		Code:       CodeForbidden,
		Message:    "forbidden",
		Errors:     errs,
		ServerTime: now(),
	})
}

func NotFound(c *fiber.Ctx, errs any) error {
	return c.Status(http.StatusNotFound).JSON(ResponseModel{
		Code:       CodeNotFound,
		Message:    "not found",
		Errors:     errs,
		ServerTime: now(),
	})
}

func Conflict(c *fiber.Ctx, errs any) error {
	return c.Status(http.StatusConflict).JSON(ResponseModel{
		Code:       CodeConflict,
		Message:    "conflict",
		Errors:     errs,
		ServerTime: now(),
	})
}

func InternalServerError(c *fiber.Ctx, errs any) error {
	return c.Status(http.StatusInternalServerError).JSON(ResponseModel{
		Code:       CodeInternalError,
		Message:    "internal server error",
		Errors:     errs,
		ServerTime: now(),
	})
}
