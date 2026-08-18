package middleware

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"backend-challenge-golang-7solution/pkg/datetime"
)

const (
	LocalKeyRequestID = "requestId"
	RequestIDHeader   = "X-Request-ID"
)

type AccessLogModel struct {
	Timestamp  string  `json:"timestamp"`
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	Status     int     `json:"status"`
	DurationMs float64 `json:"durationMs"`
	UserID     string  `json:"userId"`
	RequestID  string  `json:"requestId"`
}

func RequestID(fiberContext *fiber.Ctx) string {
	requestID, _ := fiberContext.Locals(LocalKeyRequestID).(string)
	return requestID
}

type RequestLogger interface {
	Printf(format string, values ...any)
}

func writeJSONLog(logger RequestLogger, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		logger.Printf("%s", `{"level":"error","event":"logMarshalFailed"}`)
		return
	}
	logger.Printf("%s", body)
}

func Logging(logger RequestLogger) fiber.Handler {
	if logger == nil {
		logger = log.New(os.Stderr, "", 0)
	}

	return func(fiberContext *fiber.Ctx) error {
		requestID := strings.TrimSpace(fiberContext.Get(RequestIDHeader))
		if requestID == "" {
			requestID = uuid.NewString()
		}
		fiberContext.Locals(LocalKeyRequestID, requestID)
		fiberContext.Set(RequestIDHeader, requestID)

		startedAt := time.Now()
		err := fiberContext.Next()
		if fiberContext.Path() == "/health" {
			return err
		}

		userID := UserID(fiberContext)
		if userID == "" {
			userID = "-"
		}

		writeJSONLog(logger, AccessLogModel{
			Timestamp:  datetime.FormatDateTime(startedAt),
			Method:     fiberContext.Method(),
			Path:       fiberContext.Path(),
			Status:     fiberContext.Response().StatusCode(),
			DurationMs: float64(time.Since(startedAt).Microseconds()) / 1000,
			UserID:     userID,
			RequestID:  requestID,
		})
		return err
	}
}
