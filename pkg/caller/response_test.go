package caller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type fakeLogger struct {
	called bool
	format string
	values []any
}

func (logger *fakeLogger) Printf(format string, values ...any) {
	logger.called = true
	logger.format = format
	logger.values = values
}

func TestInternalError_HidesErrorFromClientAndLogsIt(t *testing.T) {
	logger := &fakeLogger{}
	internalErrorLogger = logger
	t.Cleanup(func() {
		internalErrorLogger = nil
	})

	application := fiber.New()
	application.Get("/fail", func(fiberContext *fiber.Ctx) error {
		fiberContext.Locals(requestIDLocalKey, "req-1")
		return InternalError(fiberContext, errors.New("server selection timeout"))
	})

	request := httptest.NewRequest("GET", "/fail", nil)
	response, err := application.Test(request, -1)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	responseBodyBytes, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var responseEnvelope map[string]any
	require.NoError(t, json.Unmarshal(responseBodyBytes, &responseEnvelope))
	require.Equal(t, CodeInternalError, int(responseEnvelope["code"].(float64)))
	require.Equal(t, internalServerErrorMessage, responseEnvelope["message"])
	require.Equal(t, internalServerErrorMessage, responseEnvelope["errors"])
	require.NotContains(t, string(responseBodyBytes), "server selection timeout")

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
	require.Equal(t, "error", record["level"])
	require.Equal(t, "internalError", record["event"])
	require.Equal(t, "req-1", record["requestId"])
	require.Equal(t, "server selection timeout", record["error"])
	require.NotEmpty(t, record["timestamp"])
}
