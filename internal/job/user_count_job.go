package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	domainuser "github.com/twrnakata/user-auth-api/internal/domain/user"
	"github.com/twrnakata/user-auth-api/pkg/apperror"
	"github.com/twrnakata/user-auth-api/pkg/datetime"
)

const DefaultUserCountInterval = 10 * time.Second

type JobLogger interface {
	Printf(format string, values ...any)
}

type UserCountLogModel struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Count     int64  `json:"count"`
}

type UserCountErrorLogModel struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Event     string `json:"event"`
	Error     string `json:"error"`
	Stack     string `json:"stack,omitempty"`
}

type UserCountJob struct {
	CountUserService domainuser.CountUserService
	Logger           JobLogger
	Interval         time.Duration
}

func NewUserCountJob(countUserService domainuser.CountUserService, logger JobLogger, interval time.Duration) (*UserCountJob, error) {
	if countUserService == nil {
		return nil, apperror.ErrCountUserServiceNil
	}
	if logger == nil {
		logger = log.New(os.Stderr, "", 0)
	}
	if interval <= 0 {
		interval = DefaultUserCountInterval
	}

	return &UserCountJob{
		CountUserService: countUserService,
		Logger:           logger,
		Interval:         interval,
	}, nil
}

func (job *UserCountJob) Run(executionContext context.Context) {
	if job == nil {
		return
	}

	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	job.countAndLog(executionContext)

	for {
		select {
		case <-executionContext.Done():
			return
		case <-ticker.C:
			job.countAndLog(executionContext)
		}
	}
}

func (job *UserCountJob) countAndLog(executionContext context.Context) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		job.writeJSONLog(UserCountErrorLogModel{
			Timestamp: datetime.FormatDateTime(datetime.GetCurrentDateTimeNow()),
			Level:     "error",
			Event:     "userCountPanic",
			Error:     fmt.Sprint(recovered),
			Stack:     string(debug.Stack()),
		})
	}()

	if job.CountUserService == nil {
		job.writeJSONLog(UserCountErrorLogModel{
			Timestamp: datetime.FormatDateTime(datetime.GetCurrentDateTimeNow()),
			Level:     "error",
			Event:     "userCount",
			Error:     apperror.ErrCountUserServiceNil.Error(),
		})
		return
	}

	var count int64
	err := job.CountUserService.Count(executionContext, &count)
	if err != nil {
		job.writeJSONLog(UserCountErrorLogModel{
			Timestamp: datetime.FormatDateTime(datetime.GetCurrentDateTimeNow()),
			Level:     "error",
			Event:     "userCount",
			Error:     err.Error(),
		})
		return
	}

	job.writeJSONLog(UserCountLogModel{
		Timestamp: datetime.FormatDateTime(datetime.GetCurrentDateTimeNow()),
		Event:     "userCount",
		Count:     count,
	})
}

func (job *UserCountJob) writeJSONLog(payload any) {
	if job.Logger == nil {
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		job.Logger.Printf("%s", `{"level":"error","event":"logMarshalFailed"}`)
		return
	}
	job.Logger.Printf("%s", body)
}
