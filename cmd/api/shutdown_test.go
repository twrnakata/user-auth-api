package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type recordingHTTPServer struct {
	order           *[]string
	receivedTimeout time.Duration
	err             error
}

func (server *recordingHTTPServer) ShutdownWithTimeout(timeout time.Duration) error {
	*server.order = append(*server.order, "http")
	server.receivedTimeout = timeout
	return server.err
}

type recordingMongoClient struct {
	order *[]string
	err   error
}

func (client *recordingMongoClient) Disconnect(executionContext context.Context) error {
	*client.order = append(*client.order, "mongo")
	return client.err
}

func TestGracefulShutdown_StopsHTTPThenJobThenMongo(t *testing.T) {
	var order []string
	jobDone := &sync.WaitGroup{}
	jobDone.Add(1)

	httpServer := &recordingHTTPServer{order: &order}
	database := &recordingMongoClient{order: &order}

	err := gracefulShutdown(
		context.Background(),
		httpServer,
		func() {
			order = append(order, "cancel")
			jobDone.Done()
		},
		jobDone,
		database,
		nil,
	)

	require.NoError(t, err)
	require.Equal(t, []string{"http", "cancel", "mongo"}, order)
	require.Equal(t, httpShutdownTimeout, httpServer.receivedTimeout)
}

func TestGracefulShutdown_HTTPError_StillCancelsJobAndDisconnects(t *testing.T) {
	var order []string
	jobDone := &sync.WaitGroup{}
	jobDone.Add(1)
	httpErr := errors.New("shutdown timeout")

	err := gracefulShutdown(
		context.Background(),
		&recordingHTTPServer{order: &order, err: httpErr},
		func() {
			order = append(order, "cancel")
			jobDone.Done()
		},
		jobDone,
		&recordingMongoClient{order: &order},
		nil,
	)

	require.ErrorIs(t, err, httpErr)
	require.Equal(t, []string{"http", "cancel", "mongo"}, order)
}

func TestGracefulShutdown_DisconnectError_ReturnedAfterJobStops(t *testing.T) {
	var order []string
	jobDone := &sync.WaitGroup{}
	jobDone.Add(1)
	disconnectErr := errors.New("disconnect failed")

	err := gracefulShutdown(
		context.Background(),
		&recordingHTTPServer{order: &order},
		func() {
			order = append(order, "cancel")
			jobDone.Done()
		},
		jobDone,
		&recordingMongoClient{order: &order, err: disconnectErr},
		nil,
	)

	require.ErrorIs(t, err, disconnectErr)
	require.Equal(t, []string{"http", "cancel", "mongo"}, order)
}

type recordingShutdownLogger struct {
	lines []string
}

func (logger *recordingShutdownLogger) Printf(format string, values ...any) {
	logger.lines = append(logger.lines, format)
}

func TestGracefulShutdown_LogsHTTPThenJobThenMongo(t *testing.T) {
	jobDone := &sync.WaitGroup{}
	jobDone.Add(1)
	logger := &recordingShutdownLogger{}
	var order []string

	err := gracefulShutdown(
		context.Background(),
		&recordingHTTPServer{order: &order},
		func() {
			jobDone.Done()
		},
		jobDone,
		&recordingMongoClient{order: &order},
		logger,
	)

	require.NoError(t, err)
	require.Equal(t, []string{
		"stopping HTTP (at most %s for in-flight requests; returns now if none)",
		"HTTP stopped",
		"stopping user count job",
		"user count job stopped",
		"disconnecting Mongo",
		"Mongo disconnected",
	}, logger.lines)
}
