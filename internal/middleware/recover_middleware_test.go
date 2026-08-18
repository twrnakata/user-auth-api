package middleware

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/twrnakata/user-auth-api/pkg/caller"
)

func TestRecover_Panic_Returns500AndKeepsAppAlive(t *testing.T) {
	logger := &fakeRequestLogger{}
	application := fiber.New()
	application.Use(Recover(logger))
	application.Get("/panic", func(fiberContext *fiber.Ctx) error {
		panic("boom")
	})
	application.Get("/health", func(fiberContext *fiber.Ctx) error {
		return fiberContext.SendStatus(fiber.StatusOK)
	})

	panicRequest := httptest.NewRequest("GET", "/panic", nil)
	panicResponse, err := application.Test(panicRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusInternalServerError, panicResponse.StatusCode)

	responseBodyBytes, _ := io.ReadAll(panicResponse.Body)
	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, caller.CodeInternalError, int(responseEnvelope["code"].(float64)))
	require.Equal(t, "internal server error", responseEnvelope["errors"])
	require.NotContains(t, string(responseBodyBytes), "boom")

	record := parseLoggedJSON(t, logger)
	require.Equal(t, "error", record["level"])
	require.Equal(t, "panic", record["event"])
	require.Equal(t, "boom", record["error"])
	require.NotEmpty(t, record["stack"])

	healthRequest := httptest.NewRequest("GET", "/health", nil)
	healthResponse, err := application.Test(healthRequest, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, healthResponse.StatusCode)
}

func TestLogging_Recover_Panic_StillWritesAccessLog(t *testing.T) {
	accessLogger := &fakeRequestLogger{}
	panicLogger := &fakeRequestLogger{}
	application := fiber.New()
	application.Use(Logging(accessLogger))
	application.Use(Recover(panicLogger))
	application.Get("/panic", func(fiberContext *fiber.Ctx) error {
		panic("boom")
	})

	request := httptest.NewRequest("GET", "/panic", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	panicRecord := parseLoggedJSON(t, panicLogger)
	accessRecord := parseLoggedJSON(t, accessLogger)
	require.Equal(t, "panic", panicRecord["event"])
	require.Equal(t, float64(fiber.StatusInternalServerError), accessRecord["status"])
	require.Equal(t, panicRecord["requestId"], accessRecord["requestId"])
}
