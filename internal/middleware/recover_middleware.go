package middleware

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"

	"backend-challenge-golang-7solution/pkg/caller"
	"backend-challenge-golang-7solution/pkg/datetime"
)

type PanicLogModel struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Event     string `json:"event"`
	RequestID string `json:"requestId"`
	Error     string `json:"error"`
	Stack     string `json:"stack"`
}

func Recover(logger RequestLogger) fiber.Handler {
	if logger == nil {
		logger = log.New(os.Stderr, "", 0)
	}

	return func(fiberContext *fiber.Ctx) error {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			requestID := RequestID(fiberContext)
			if requestID == "" {
				requestID = "-"
			}

			writeJSONLog(logger, PanicLogModel{
				Timestamp: datetime.FormatDateTime(datetime.GetCurrentDateTimeNow()),
				Level:     "error",
				Event:     "panic",
				RequestID: requestID,
				Error:     fmt.Sprint(recovered),
				Stack:     string(debug.Stack()),
			})
			_ = caller.InternalServerError(fiberContext, "internal server error")
		}()

		return fiberContext.Next()
	}
}
