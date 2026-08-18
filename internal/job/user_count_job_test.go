package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeCountUserService struct {
	mu     sync.Mutex
	called int
	count  int64
	err    error
}

func (service *fakeCountUserService) Count(executionContext context.Context, count *int64) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.called++
	if count != nil {
		*count = service.count
	}
	return service.err
}

func (service *fakeCountUserService) callCount() int {
	service.mu.Lock()
	defer service.mu.Unlock()
	return service.called
}

type recordingLogger struct {
	mu    sync.Mutex
	lines []string
}

func (logger *recordingLogger) Printf(format string, values ...any) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.lines = append(logger.lines, fmt.Sprintf(format, values...))
}

func (logger *recordingLogger) snapshot() []string {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	copied := make([]string, len(logger.lines))
	copy(copied, logger.lines)
	return copied
}

type panickingCountUserService struct{}

func (service *panickingCountUserService) Count(executionContext context.Context, count *int64) error {
	panic("count boom")
}

func TestNewUserCountJob_NilService_ReturnsError(t *testing.T) {
	job, err := NewUserCountJob(nil, &recordingLogger{}, time.Second)
	require.Error(t, err)
	require.Nil(t, job)
}

func TestUserCountJob_CountAndLog_WritesUserCountJSON(t *testing.T) {
	logger := &recordingLogger{}
	service := &fakeCountUserService{count: 3}
	job, err := NewUserCountJob(service, logger, time.Second)
	require.NoError(t, err)

	job.countAndLog(context.Background())

	lines := logger.snapshot()
	require.Len(t, lines, 1)

	var payload UserCountLogModel
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &payload))
	require.Equal(t, "userCount", payload.Event)
	require.Equal(t, int64(3), payload.Count)
	require.NotEmpty(t, payload.Timestamp)
}

func TestUserCountJob_CountAndLog_ServiceError_WritesErrorJSON(t *testing.T) {
	logger := &recordingLogger{}
	service := &fakeCountUserService{err: errors.New("network timeout")}
	job, err := NewUserCountJob(service, logger, time.Second)
	require.NoError(t, err)

	job.countAndLog(context.Background())

	lines := logger.snapshot()
	require.Len(t, lines, 1)

	var payload UserCountErrorLogModel
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &payload))
	require.Equal(t, "error", payload.Level)
	require.Equal(t, "userCount", payload.Event)
	require.Equal(t, "network timeout", payload.Error)
}

func TestUserCountJob_CountAndLog_Panic_DoesNotCrash(t *testing.T) {
	logger := &recordingLogger{}
	job, err := NewUserCountJob(&panickingCountUserService{}, logger, time.Second)
	require.NoError(t, err)

	require.NotPanics(t, func() {
		job.countAndLog(context.Background())
	})

	lines := logger.snapshot()
	require.Len(t, lines, 1)

	var payload UserCountErrorLogModel
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &payload))
	require.Equal(t, "error", payload.Level)
	require.Equal(t, "userCountPanic", payload.Event)
	require.Contains(t, payload.Error, "count boom")
	require.NotEmpty(t, payload.Stack)
}

func TestUserCountJob_Run_LogsUntilContextCanceled(t *testing.T) {
	logger := &recordingLogger{}
	service := &fakeCountUserService{count: 3}
	job, err := NewUserCountJob(service, logger, 15*time.Millisecond)
	require.NoError(t, err)

	executionContext, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		job.Run(executionContext)
		close(done)
	}()

	require.Eventually(t, func() bool {
		return service.callCount() >= 2
	}, time.Second, 5*time.Millisecond)

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("job did not stop after context cancel")
	}

	lines := logger.snapshot()
	require.GreaterOrEqual(t, len(lines), 2)

	var payload UserCountLogModel
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &payload))
	require.Equal(t, "userCount", payload.Event)
	require.Equal(t, int64(3), payload.Count)
}
