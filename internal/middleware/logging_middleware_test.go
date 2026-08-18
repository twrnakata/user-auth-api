package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"backend-challenge-golang-7solution/pkg/datetime"
)

type fakeRequestLogger struct {
	called bool
	format string
	values []any
}

func (logger *fakeRequestLogger) Printf(format string, values ...any) {
	logger.called = true
	logger.format = format
	logger.values = values
}

func parseLoggedJSON(t *testing.T, logger *fakeRequestLogger) map[string]any {
	t.Helper()
	require.True(t, logger.called)
	require.Len(t, logger.values, 1)

	payload, ok := logger.values[0].([]byte)
	if !ok {
		asString, stringOK := logger.values[0].(string)
		require.True(t, stringOK)
		payload = []byte(asString)
	}

	var record map[string]any
	require.NoError(t, json.Unmarshal(payload, &record))
	return record
}

func TestLogging_Success_LogsMethodPathStatusAndDuration(t *testing.T) {
	logger := &fakeRequestLogger{}
	application := fiber.New()
	application.Use(Logging(logger))
	application.Post("/auth/login", func(fiberContext *fiber.Ctx) error {
		return fiberContext.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest("POST", "/auth/login", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)

	record := parseLoggedJSON(t, logger)
	_, parseErr := time.Parse(datetime.DefaultDateTimeFormat, record["timestamp"].(string))
	require.NoError(t, parseErr)
	require.Equal(t, "POST", record["method"])
	require.Equal(t, "/auth/login", record["path"])
	require.Equal(t, float64(fiber.StatusOK), record["status"])
	require.GreaterOrEqual(t, record["durationMs"].(float64), 0.0)
	require.Equal(t, "-", record["userId"])
	require.NotEmpty(t, record["requestId"])
	require.Equal(t, record["requestId"], response.Header.Get(RequestIDHeader))
}

func TestLogging_Health_SkipsLog(t *testing.T) {
	logger := &fakeRequestLogger{}
	application := fiber.New()
	application.Use(Logging(logger))
	application.Get("/health", func(fiberContext *fiber.Ctx) error {
		return fiberContext.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest("GET", "/health", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)
	require.False(t, logger.called)
	require.NotEmpty(t, response.Header.Get(RequestIDHeader))
}

func TestLogging_NotFound_LogsStatus(t *testing.T) {
	logger := &fakeRequestLogger{}
	application := fiber.New()
	application.Use(Logging(logger))
	application.Get("/users/:id", func(fiberContext *fiber.Ctx) error {
		return fiberContext.SendStatus(fiber.StatusNotFound)
	})

	request := httptest.NewRequest("GET", "/users/missing", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, response.StatusCode)

	record := parseLoggedJSON(t, logger)
	require.Equal(t, "GET", record["method"])
	require.Equal(t, "/users/missing", record["path"])
	require.Equal(t, float64(fiber.StatusNotFound), record["status"])
}

func TestLogging_Authenticated_LogsUserID(t *testing.T) {
	logger := &fakeRequestLogger{}
	application := fiber.New()
	application.Use(Logging(logger))
	application.Get("/users", func(fiberContext *fiber.Ctx) error {
		fiberContext.Locals(LocalKeyUserID, "u-1")
		return fiberContext.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest("GET", "/users", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)
	require.Equal(t, "u-1", parseLoggedJSON(t, logger)["userId"])
}

func TestLogging_IncomingRequestID_ReusesHeader(t *testing.T) {
	logger := &fakeRequestLogger{}
	application := fiber.New()
	application.Use(Logging(logger))
	application.Get("/users", func(fiberContext *fiber.Ctx) error {
		return fiberContext.SendStatus(fiber.StatusOK)
	})

	request := httptest.NewRequest("GET", "/users", nil)
	request.Header.Set(RequestIDHeader, "req-from-client")
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, response.StatusCode)
	require.Equal(t, "req-from-client", response.Header.Get(RequestIDHeader))
	require.Equal(t, "req-from-client", parseLoggedJSON(t, logger)["requestId"])
}
